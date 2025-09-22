package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type simpleQueueType int

const (
	DurableQueue simpleQueueType = iota
	TransientQueue
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType simpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	durable, autoDelete, exclusive := GetQueueOptions(queueType)

	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil

}

func GetQueueOptions(queueType simpleQueueType) (durable, autoDelete, exclusive bool) {
	switch queueType {
	case DurableQueue:
		return true, false, false
	case TransientQueue:
		return false, true, true
	default:
		return false, false, false
	}
}
