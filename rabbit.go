package main
import(
	"context"
	"log"
	"time"
	"fmt"
	"bufio"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)


func sendMessage(ch *amqp.Channel, ctx context.Context, msg, exchange string){
	err := ch.PublishWithContext(ctx,
		exchange, //exchange
		"", //routing Key
		false, //mandatory
		false, //immediate
		amqp.Publishing{
			ContentType: 	"text/plain",
			Body:			[]byte(msg),
		},
	)
	failOnError(err, "Failed to publish message")
	
}


func runRabbit(username, password string){
	//Set up
	//username := "guest"
	//password := "guest"
	url := "localhost"
	port := "5672"
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s",username,password,url,port))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	exchangeName := "chat"
	err = ch.ExchangeDeclare(
		exchangeName,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare exchange")
	q, err := ch.QueueDeclare(
		"", // Name
		false, //durable
		false, //delete when unused
		true, //exclusive
		false, //no-wait
		nil, //args
	)
	failOnError(err, "Failed to declare queue")
	//Consume message
	msgs, err := ch.Consume(
		q.Name, //queue
		"",		//consumer
		true,	//auto-ack
		false,	//exclusive
		false,	//no-local
		false,	//no-wait
		nil,	//args
	)
	failOnError(err, "Failed to register consumer")
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind Queue")
	var forever chan struct{}

	go func()  {
		for d := range msgs{
			log.Printf("%s",d.Body)
		}
	}()

	//Send Message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	for true{
		scanner := bufio.NewScanner(os.Stdin)
    	scanner.Scan()
		sendMessage(ch,ctx,scanner.Text(),exchangeName)
	}
	

	<-forever
}