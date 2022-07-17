package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		return &logs.LogResponse{ Result: "failed" }, err
	}

	res := &logs.LogResponse{ Result: "logged by grpc" }
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Panic("Cannot listen for grpc", err)
	}

	server := grpc.NewServer()
	logs.RegisterLogServiceServer(server, &LogServer{
		Models: app.Models,
	})

	log.Printf("Started gRPC server on port %s", gRpcPort)

	err = server.Serve(lis)
	if err != nil {
		log.Panic("Cannot serve grpc", err)
	}
}
