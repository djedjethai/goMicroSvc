package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/djedjethai/listener/pkg/consumer"
	"github.com/djedjethai/listener/pkg/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	rabbitConn, err := connection()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Listening and consuming RabbitMQ messages")

	evt := event.NewEvent()
	cons, err := consumer.NewConsumer(rabbitConn, evt)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = cons.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
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
