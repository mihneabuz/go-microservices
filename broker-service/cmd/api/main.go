package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "3001"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	log.Printf("Starting broker service on port %s\n", webPort)

	rabbit, err := connect()
	if err != nil {
		log.Panic("Could not connect to RabbitMQ")
	}
	defer rabbit.Close()

	app := Config{
		Rabbit: rabbit,
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
