import pika, sys, os, threading


connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
channel = connection.channel()
def callback(ch, method, properties, body):
    print(" [x] Received %r" % body, flush=True)

def sendmesg(msg):
    channel.basic_publish(exchange='',
                        routing_key='hello',
                        body=msg)
def main():
    channel.queue_declare(queue='hello')

    channel.basic_consume(queue='hello',
                        auto_ack=True,
                        on_message_callback=callback)
    #sendmesg("hello world")
    print(' [*] Waiting for messages. To exit press CTRL+C', flush=True)
    channel.start_consuming()

threading.Thread(main())

text = input("input")
print(text)