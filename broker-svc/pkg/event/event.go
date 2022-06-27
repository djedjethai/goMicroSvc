package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Event interface {
	DeclareExchange(ch *amqp.Channel) error
	DeclareRandomQueue(ch *amqp.Channel) (amqp.Queue, error)
}

type event struct{}

func NewEvent() Event {
	return &event{}
}

func (e *event) DeclareExchange(ch *amqp.Channel) error {
	// build in func
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type
		true,         // durable ?
		false,        // autodeleted ?
		false,        // is this an exchange which is just use internally ?
		false,        // no-wait ?
		nil,          // arguments
	)
}

func (e *event) DeclareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unuse ?
		true,  // exclusive(don't share it around)
		false, // no-wait ?
		nil,   // arguments(no specific targets)
	)
}
