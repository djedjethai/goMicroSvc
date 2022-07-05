package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/djedjethai/listener/pkg/event"
	"github.com/djedjethai/toolbox"
	amqp "github.com/rabbitmq/amqp091-go"
)

var tools = new(toolbox.Tools)

type Consumer interface {
	Listen(topics []string) error
}

type consumer struct {
	conn      *amqp.Connection
	queueName string
	event     event.Event
}

func NewConsumer(conn *amqp.Connection, evt event.Event) (Consumer, error) {
	cons := &consumer{
		conn:  conn,
		event: evt,
	}

	err := cons.setup()
	if err != nil {
		return &consumer{}, err
	}

	return cons, nil
}

// listen to the queue for specific topic
func (c *consumer) Listen(topics []string) error {
	// get from the receiver consumer
	// get the channel
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// get a random queue
	q, err := c.event.DeclareRandomQueue(ch)
	if err != nil {
		return err
	}

	// range on the topics
	for _, s := range topics {
		// bind our channel to each of this topics
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false, // wait?
			nil,   // no arguments
		)

		if err != nil {
			return err
		}
	}

	// look for message
	messages, err := ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-aknowledge
		false, // is it exclusive
		false, // internal?
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// consume all the things coming from rabbitMQ forever(until the app stop)
	forever := make(chan bool)
	go func() {
		// d is the current iteration
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			// declare a second go routine to make things as fast as possible
			go c.handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	// this will make it going forever
	<-forever

	return nil
}

func (c *consumer) setup() error {
	// channel from the receiver consumer
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return c.event.DeclareExchange(channel)
}

func (c *consumer) logEvent(entry Payload) error {

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}

func (c *consumer) handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log what ever we get
		err := c.logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// where we authenticate
	// can add as many cases as we want...
	default:
		err := c.logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}
