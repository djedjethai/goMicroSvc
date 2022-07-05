package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/djedjethai/broker/pkg/consumer"
	"github.com/djedjethai/broker/pkg/domain/jsdomain"
	"github.com/djedjethai/broker/pkg/emitter"
	"github.com/djedjethai/broker/pkg/event"
	"github.com/djedjethai/broker/pkg/http/rest"
	"github.com/djedjethai/broker/pkg/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

// type Config struct{}

func main() {

	// set the connection to rabbitMQ
	rabbitConn, err := connection()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	evt := event.NewEvent()
	// here the brocker simply emit event(so the consumer is useless)
	cons, err := consumer.NewConsumer(rabbitConn, evt)
	emit, err := emitter.NewEmitter(rabbitConn, evt)
	if err != nil {
		panic(err)
	}

	// app := Config{}

	log.Println("Starting broker service on port: ", webPort)

	// get the memory datas(need to create them)
	var svc service.Service

	repo := new(jsdomain.Repository)
	// HERE the cons is useless, BUT I could pass it like so
	svc = service.NewService(repo, cons, emit)

	// populate the repo
	svc.AddData(sampleData)

	mux := rest.Handler(svc)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: mux,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connection() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			fmt.Println("RabbitMQ not yet ready")
			counts++
		} else {
			log.Println("Connected to rabbitMQ!")
			connection = c
			break
		}

		// if we can not connect, we don't want to run the try endlessly
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// we want to iincrease the delay each time we backOff
		// raise its duration to power of two
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
