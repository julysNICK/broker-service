package event

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer is for creating a new Consumer (it returns a Consumer object)

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// func setup() error is for setting up the Consumer (it creates a channel and declares the exchange)
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchangePost(channel)
}

// Payload is for storing the payload of a message (it will be used for unmarshalling the message)
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// func Listen(topics []string) error is for listening to messages published to the exchange (it creates a queue, binds it to the exchange, and consumes messages from the queue) (it returns an error if any)
func (consumer *Consumer) ListenPost(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueuePost(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"events",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePostPayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func (consumer *Consumer) ListenUser(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueuePost(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"events_user",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePostPayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

// func handlePayload(payload Payload) is for handling the payload of a message (it prints the payload to the console)
func handlePostPayload(payload Payload) {
	fmt.Printf("Received message: %v\n", payload)

	switch payload.Name {
	case "post.created":
		fmt.Println("post.created")
	case "post.deleted":
		fmt.Println("post.deleted")
	case "post.updated":
		fmt.Println("post.updated")

	default:
		fmt.Println("unknown event")

	}
}
