package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"qq/pkg/rabbitqq"
	"qq/services/qq"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ThreadCount = 20

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

	var wg sync.WaitGroup

	for i := 0; i < ThreadCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for msg := range msgs {
				ctx := context.Background()
				err = s.handleRawMessage(ctx, msg.Body, msg.CorrelationId, msg.ReplyTo)
				if err != nil {
					log.Println(fmt.Errorf("failed to handle message: %w", err))
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func handleMessage[Message any, ReplyMessage any](ctx context.Context, s server, body []byte, corrId string, replyTo string, proc func(message Message) ReplyMessage) error {
	var message Message
	err := json.Unmarshal(body, &message)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	replyMessage := proc(message)

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

func (s server) handleRawMessage(ctx context.Context, body []byte, corrId string, replyTo string) error {
	var baseMessage rabbitqq.BaseMessage
	err := json.Unmarshal(body, &baseMessage)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	switch baseMessage.Name {
	case "add":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(addMessage rabbitqq.AddMessage) rabbitqq.AddReplyMessage {
				entity := FromAddMessage(addMessage)
				return ToAddReplyMessage(s.service.Add(ctx, entity))
			})
		if err != nil {
			return fmt.Errorf("failed to handle add message: %w", err)
		}

	case "remove":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(removeMessage rabbitqq.RemoveMessage) rabbitqq.RemoveReplyMessage {
				key := FromRemoveMessage(removeMessage)
				return ToRemoveReplyMessage(s.service.Remove(ctx, key))
			})
		if err != nil {
			return fmt.Errorf("failed to handle remove message: %w", err)
		}

	case "get":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(getMessage rabbitqq.GetMessage) rabbitqq.GetReplyMessage {
				key := FromGetMessage(getMessage)
				return ToGetReplyMessage(s.service.Get(ctx, key))
			})
		if err != nil {
			return fmt.Errorf("failed to handle get message: %w", err)
		}

	case "get all":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(getAllMessage rabbitqq.GetAllMessage) rabbitqq.GetAllReplyMessage {
				return ToGetAllReplyMessage(s.service.GetAll(ctx))
			})
		if err != nil {
			return fmt.Errorf("failed to handle add message: %w", err)
		}
	}

	return nil
}
