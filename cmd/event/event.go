package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchangePost(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable?
		false,        // auto-deleted?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguements?
	)
}

func declareExchangeUser(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"events_user", // name
		"topic",       // type
		true,          // durable?
		false,         // auto-deleted?
		false,         // internal?
		false,         // no-wait?
		nil,           // arguements?
	)
}

func declareRandomQueuePost(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unused?
		true,  // exclusive?
		false, // no-wait?
		nil,   // arguments?
	)
}
