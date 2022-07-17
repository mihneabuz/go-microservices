package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

type RPCServer struct {}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	log.Println("RPC LogInfo")

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	*resp = "logged with rpc"
	return nil
}
