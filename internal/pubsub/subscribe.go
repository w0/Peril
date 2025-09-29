package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
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
				log.Printf("error unmarshaling delivery msg: %w", err)
			}

			handler(parsed)

			msg.Ack(false)
		}
	}(deliveries)

	return nil
}
