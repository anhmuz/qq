package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"qq/pkg/rabbitqq"
	"qq/services/qq"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server interface {
	Serve() error
}

type server struct {
	queue   string
	service qq.Service
	channel *amqp.Channel
}

var _ Server = server{}

func NewServer(queue string, service qq.Service) (Server, error) {
	log.Printf("create new RabbitMQ server - queue: %v\n ", queue)

	ch, err := connect(queue)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &server{
		queue:   queue,
		service: service,
		channel: ch,
	}, nil
}

func connect(queue string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(rabbitqq.AmqpServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return ch, nil
}

func (s server) Serve() error {
	msgs, err := s.channel.Consume(
		s.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	for msg := range msgs {
		var message interface{}
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		s.parse(message, msg.CorrelationId, msg.ReplyTo)
	}

	return nil
}

func (s server) parse(message interface{}, corrId string, replyTo string) {
	convertor, _ := NewConvertor()

	switch message.(type) {
	case rabbitqq.AddMessage:
		entity := convertor.AddMessageToEntity(message.(rabbitqq.AddMessage))
		addReplyMessage := convertor.BoolToAddReplyMessage(s.service.Add(entity))
		s.sendMessage(addReplyMessage, corrId, replyTo)

	case rabbitqq.RemoveMessage:
		key := convertor.RemoveMessageToString(message.(rabbitqq.RemoveMessage))
		removeReplyMessage := convertor.BoolToRemoveReplyMessage(s.service.Remove(key))
		s.sendMessage(removeReplyMessage, corrId, replyTo)

	case rabbitqq.GetMessage:
		key := convertor.GetMessageToString(message.(rabbitqq.GetMessage))
		getReplyMessage := convertor.EntityToGetReplyMessage(s.service.Get(key))
		s.sendMessage(getReplyMessage, corrId, replyTo)

	case rabbitqq.GetAllMessage:
		getAllReplyMessage := convertor.EntitiesToGetAllReplyMessage(s.service.GetAll())
		s.sendMessage(getAllReplyMessage, corrId, replyTo)
	}
}

func (s server) sendMessage(message interface{}, corrId string, replyTo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to produce JSON: %w", err)
	}

	err = s.channel.PublishWithContext(ctx,
		"",
		replyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			Body:          jsonMessage,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
