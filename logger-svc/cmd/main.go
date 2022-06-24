package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/djedjethai/logger/pkg/handlers/handlehttp"
	"github.com/djedjethai/logger/pkg/service"
	"github.com/djedjethai/logger/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Type int

const (
	webPort       = "80"
	rpcPort       = "5001"
	mongoURL      = "mongodb://mongo:27017"
	gRpcPort      = "50001"
	HTTP     Type = iota
)

// var client *mongo.Client

func main() {

	protocole := HTTP

	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	// client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// create store
	store := storage.New(mongoClient)

	// create svc
	service := service.NewService(store.LogEntry)

	// get the handler
	mux := handlehttp.Handler(service)

	switch protocole {
	case HTTP:
		// go serveHTTP(mux)
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", webPort),
			Handler: mux,
		}

		log.Println("Logger listen on port ", webPort)
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}

	}

	// close connection
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// func serveHTTP(mux http.Handler) {
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webPort),
// 		Handler: mux,
// 	}
//
// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

func connectToMongo() (*mongo.Client, error) {
	//  create a connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting: ", err)
		return nil, err
	}

	return c, nil
}
