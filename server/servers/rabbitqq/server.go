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
		err = s.handleRawMessage(msg.Body, msg.CorrelationId, msg.ReplyTo)
		if err != nil {
			return fmt.Errorf("failed to handle message: %w", err)
		}
	}

	return nil
}

func handleMessage[Message any, ReplyMessage any](s server, body []byte, corrId string, replyTo string, proc func(message Message) ReplyMessage) error {
	var message Message
	err := json.Unmarshal(body, &message)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	replyMessage := proc(message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonReplyMessage, err := json.Marshal(replyMessage)
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
			Body:          jsonReplyMessage,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a reply message: %w", err)
	}

	return nil
}

func (s server) handleRawMessage(body []byte, corrId string, replyTo string) error {
	var baseMessage rabbitqq.BaseMessage
	err := json.Unmarshal(body, &baseMessage)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	convertor, _ := NewConvertor()

	switch baseMessage.Name {
	case "add":
		err = handleMessage(s, body, corrId, replyTo,
			func(addMessage rabbitqq.AddMessage) rabbitqq.AddReplyMessage {
				entity := convertor.AddMessageToEntity(addMessage)
				a := convertor.BoolToAddReplyMessage(s.service.Add(entity))
				return a
			})
		if err != nil {
			return fmt.Errorf("failed to handle add message: %w", err)
		}

	case "remove":
		err = handleMessage(s, body, corrId, replyTo,
			func(removeMessage rabbitqq.RemoveMessage) rabbitqq.RemoveReplyMessage {
				key := convertor.RemoveMessageToString(removeMessage)
				return convertor.BoolToRemoveReplyMessage(s.service.Remove(key))
			})
		if err != nil {
			return fmt.Errorf("failed to handle remove message: %w", err)
		}

	case "get":
		err = handleMessage(s, body, corrId, replyTo,
			func(getMessage rabbitqq.GetMessage) rabbitqq.GetReplyMessage {
				key := convertor.GetMessageToString(getMessage)
				return convertor.EntityToGetReplyMessage(s.service.Get(key))
			})
		if err != nil {
			return fmt.Errorf("failed to handle get message: %w", err)
		}

	case "get all":
		err = handleMessage(s, body, corrId, replyTo,
			func(getAllMessage rabbitqq.GetAllMessage) rabbitqq.GetAllReplyMessage {
				return convertor.EntitiesToGetAllReplyMessage(s.service.GetAll())
			})
		if err != nil {
			return fmt.Errorf("failed to handle add message: %w", err)
		}
	}

	return nil
}
