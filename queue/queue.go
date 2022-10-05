package queue

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reciever-ms/config"
	"reciever-ms/queue/interfaces"
	"reciever-ms/tracer"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/codes"
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
	go func() {
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
	}()
}

func (s *Server) PublishTrackerEvent(ctx context.Context, evt TrackerEvent) error {
	ctx, span := tracer.NewSpan(ctx, "queue", "PublishTrackerEvent")
	defer span.End()

	if evt.Protocol == "" || evt.Type == "" || evt.Imei == "" {
		errMsg := "cannot mount routing key: protocol, type and imei must not be empty"
		span.SetStatus(codes.Error, errMsg)

		return errors.New(errMsg)
	}

	body, err := json.Marshal(evt)
	if err != nil {
		tracer.AddSpanErrorAndFail(span, err, "failed to marshal tracker event")
		return err
	}

	routingKey := strings.Join([]string{evt.Protocol, evt.Type, evt.Imei}, ".")

	err = s.Publisher.PublishWithContext(ctx, s.channel, s.cfg.Exchange, routingKey, amqp.Publishing{Body: body})
	if err != nil {
		tracer.AddSpanErrorAndFail(span, err, "failed to publish tracker event")
	}

	return err
}

func (s *Server) Stop() error {
	log.Printf("[RMQ] closing connections")

	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}
