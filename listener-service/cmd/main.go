package main

import (
	"listener/event"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting listerer service")

	rabbit, err := connect()
	if err != nil {
		log.Panic("Could not connect to RabbitMQ")
	}
	defer rabbit.Close()

	log.Println("Listening and consuming RabbitMQ")

	consumer, err := event.NewComsumer(rabbit)
	if err != nil {
		log.Panic("Could not consume from RabbitMQ")
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic("Could not consume from RabbitMQ")
	}

}

func connect() (*amqp.Connection, error) {
	tries, maxTries := 0, 16

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("Cannot connect to RabbitMQ")
			tries += 1
			time.Sleep(2 * time.Second)
		} else {
			return c, nil
		}

		if tries >= maxTries {
			return nil, err
		}
	}
}
