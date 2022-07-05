package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type RPCServer struct {
	clt *mongo.Client
}

type RPCPayload struct {
	Name string
	Data string
}

func NewRPCServer(c *mongo.Client) *RPCServer {
	return &RPCServer{
		clt: c,
	}
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := r.clt.Database("logs").Collection("logs")

	fmt.Println("in RPC LogInfo")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs: ", err)
	}

	*resp = "Processed payload via RPC:" + payload.Name

	return nil
}
