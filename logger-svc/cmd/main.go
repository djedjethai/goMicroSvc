package main

// mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&appname
// =MongoDB%20Compass&directConnection=true&ssl=false

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/djedjethai/logger/pkg/handlers/handlehttp"
	// "github.com/djedjethai/logger/pkg/handlers/handlerpc"
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
	RPC
)

// var client *mongo.Client

func main() {

	// protocole := HTTP
	protocole := RPC

	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	switch protocole {
	case HTTP:
		// create store
		store := storage.New(mongoClient)

		// create svc
		service := service.NewService(store.LogEntry)

		// get the handler
		mux := handlehttp.Handler(service)

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
	case RPC:
		fmt.Println("in RPC main")

		// pass a pointer to the mongo client to be use there
		rpcServer := storage.NewRPCServer(mongoClient)

		// Register the rpc server(tell the app that the rpc server listen)
		err := rpc.Register(rpcServer)
		if err != nil {
			log.Panic(err)
		}

		// ========= Need to start a server anyway ==========
		// go serveHTTP(mux)
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", webPort),
			Handler: nil,
		}

		// start rpc server
		go rpcListen()

		log.Println("Logger listen on port ", webPort)
		err = srv.ListenAndServe()
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

func rpcListen() error {
	log.Println("Starting RPC server on port: ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		fmt.Println("see err logger rpc connecting: ", err)
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

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
