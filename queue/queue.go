package queue

import (
	"log"
	"reciever-ms/config"
	"reciever-ms/queue/interfaces"

	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate mockgen -destination=../mocks/amqp.go -package=mocks github.com/rabbitmq/amqp091-go Acknowledger

type Server struct {
	interfaces.Connector
	interfaces.Publisher

	cfg         config.RmqConfig
	conn        interfaces.AmqpConnection
	channel     interfaces.AmqpChannel
	notifyClose chan *amqp.Error
}

func New(cfg config.RmqConfig) Server {
	return Server{cfg: cfg, Connector: &Connector{}, Publisher: &Publisher{}}
}

func (s *Server) Start() {
	for {
		s.connect()

		connectionError, chanClosed := <-s.notifyClose

		// connection error is nil and chanClosed is false when
		// the connection was closed manually with client code
		if connectionError != nil {
			log.Printf("[RMQ] connection error: %v \n", connectionError)
		}

		if !chanClosed {
			return
		}
	}
}

func (s *Server) Stop() error {
	log.Printf("[RMQ] closing connections")

	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}
