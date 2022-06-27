package emitter

import (
	"log"

	"github.com/djedjethai/broker/pkg/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter interface {
	Push(event string, severity string) error
}

type emitter struct {
	conn *amqp.Connection
	evt  event.Event
}

func NewEmitter(c *amqp.Connection, e event.Event) (Emitter, error) {
	emit := &emitter{
		conn: c,
		evt:  e,
	}

	err := emit.setup()
	if err != nil {
		return &emitter{}, err
	}

	return emit, nil
}

func (e *emitter) Push(event string, severity string) error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	// call the channel's build-in method Publish
	err = channel.Publish(
		"logs_topic",
		severity, // ehat key we are using, we have 3 log.INFO, log.WARNING, etc
		false,    // mandatory
		false,    // is it immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event), // cast to a slice of bytes event
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (e *emitter) setup() error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return e.evt.DeclareExchange(channel)
}
