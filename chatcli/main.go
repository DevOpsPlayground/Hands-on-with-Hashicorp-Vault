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
func failOnError(err error, msg string){
	if err != nil{
		log.Panicf("%s: %s", msg, err)
	}
}

func sendMessage(ch *amqp.Channel, ctx context.Context, msg, queueName string){
	err := ch.PublishWithContext(ctx,
		"", //exchange
		queueName, //routing Key
		false, //mandatory
		false, //immediate
		amqp.Publishing{
			ContentType: 	"text/plain",
			Body:			[]byte(msg),
		},
	)
	failOnError(err, "Failed to publish message")
	
}

func main(){
	//Set up
	username := "guest"
	password := "guest"
	url := "localhost"
	port := "5672"
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s",username,password,url,port))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // Name
		false, //durable
		false, //delete when unused
		false, //exclusive
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
		sendMessage(ch,ctx,scanner.Text(),q.Name)
	}
	

	<-forever
}