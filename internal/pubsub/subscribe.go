package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {

	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(d <-chan amqp.Delivery) {
		for msg := range d {
			var parsed T
			err := json.Unmarshal(msg.Body, &parsed)
			if err != nil {
				log.Printf("error unmarshaling delivery msg: %s\n", err)
			}

			ackAction := handler(parsed)

			switch ackAction {
			case Ack:
				msg.Ack(false)
				log.Println("msg Ack.")
			case NackRequeue:
				msg.Nack(false, true)
				log.Println("msg NackRequeue.")
			case NackDiscard:
				msg.Nack(false, false)
				log.Println("msg NackDiscard.")
			}
		}
	}(deliveries)

	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(d <-chan amqp.Delivery) {
		for msg := range d {
			var parsed T
			buf := bytes.NewBuffer(msg.Body)

			dec := gob.NewDecoder(buf)
			err := dec.Decode(&parsed)
			if err != nil {
				log.Printf("error decoding game_log: %s\n", err)
			}

			ackAction := handler(parsed)

			switch ackAction {
			case Ack:
				msg.Ack(false)
				log.Println("msg Ack.")
			case NackRequeue:
				msg.Nack(false, true)
				log.Println("msg NackRequeue.")
			case NackDiscard:
				msg.Nack(false, false)
				log.Println("msg NackDiscard.")
			}
		}

	}(deliveries)

	return nil
}
