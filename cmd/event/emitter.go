package event

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter is for storing the connection to RabbitMQ (we'll use it to create a channel)
type Emitter struct {
	connection *amqp.Connection
}

// func (e *Emitter) setup() error is for setting up the Emitter (it creates a channel and declares the exchange)
func (e *Emitter) setupRabbitPost() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	return declareExchangePost(channel)
}

func (e *Emitter) setupRabbitUser() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	return declareExchangeUser(channel)
}

// func (e *Emitter) Push(event string, severity string) error is for publishing a message to the exchange (it returns an error if any) (it will be used for publishing messages)
func (e *Emitter) PushPost(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		fmt.Println("line 29-emitter.go " + err.Error())
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")
	log.Println("event: " + event)
	log.Println("severity: " + severity)
	// channel.PublishWithContext is for publishing a message to the exchange (it returns an error if any)
	err = channel.PublishWithContext(
		context.Background(),
		"events",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		fmt.Println("line 49-emitter.go " + err.Error())
		return err
	}

	return nil
}

func (e *Emitter) PushUser(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		fmt.Println("line 29-emitter.go " + err.Error())
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")
	log.Println("event: " + event)
	log.Println("severity: " + severity)
	// channel.PublishWithContext is for publishing a message to the exchange (it returns an error if any)
	err = channel.PublishWithContext(
		context.Background(),
		"events_user",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		fmt.Println("line 49-emitter.go " + err.Error())
		return err
	}

	return nil
}

//func NewEventEmitter(conn *amqp.Connection) (Emitter, error) is for creating a new Emitter (it returns an Emitter object) (it will be used for creating a new Emitter)
//Emitter is for storing the connection to RabbitMQ (we'll use it to create a channel)

func NewEventEmitterPost(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setupRabbitPost()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}

func NewEventEmitterUser(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setupRabbitUser()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
