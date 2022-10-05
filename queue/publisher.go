package queue

import (
	"context"
	"reciever-ms/queue/interfaces"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct{}

func (p *Publisher) PublishWithContext(channel interfaces.AmqpChannel, ctx context.Context, exchange, key string, msg amqp.Publishing) error {
	return channel.PublishWithContext(
		ctx,      // context
		exchange, // exchange
		key,      // key
		false,    // mandatory
		false,    // immediate
		msg,      // msg
	)
}
