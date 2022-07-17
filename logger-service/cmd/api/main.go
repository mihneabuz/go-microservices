package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	rpcPort = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"

	maxTries = 16
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting logger service on port", webPort)

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	} ()

	app := Config{
		Models: data.New(client),
	}

	rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	app.serve()
}

func (app *Config) serve() {
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() {
	log.Println("Listening to rpc on port " + rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Panic("Cannot listen to rpc")
	}
	defer listen.Close()

	for {
		rpcCon, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcCon)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "mongo",
		Password: "mongo",
	})

	tries := 0
	for {
		c, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Println("Error connecting to mongo:", err)
			tries += 1
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Connected to mongo")
			return c, nil
		}

		if tries >= maxTries {
			return nil, err
		}
	}
}

