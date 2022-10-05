package interfaces

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate mockgen -destination=../../mocks/server.go -package=mocks . IServer,AmqpChannel,AmqpConnection,Connector,Publisher

type AmqpChannel interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

type AmqpDialer func(url string) (AmqpConnection, error)

type AmqpConnection interface {
	Close() error
	Channel() (AmqpChannel, error)
}

type Connector interface {
	Connect(url string) (AmqpConnection, error)
}

type Publisher interface {
	PublishWithContext(channel AmqpChannel, ctx context.Context, exchange, key string, msg amqp.Publishing) error
}

type IServer interface {
	Connector
	Publisher

	ServeDelivery(d *amqp.Delivery)
	Shutdown() error
}
