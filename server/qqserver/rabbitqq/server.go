package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"qq/pkg/log"
	"qq/pkg/protocol"
	"qq/pkg/qqclient/rabbitqq"
	"qq/pkg/qqcontext"
	"qq/server/qqserver"
	"qq/services/qq"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ThreadCount = 20

type server struct {
	queue   string
	service qq.Service
	channel *amqp.Channel
}

var _ qqserver.Server = server{}

func NewServer(ctx context.Context, queue string, service qq.Service) (qqserver.Server, error) {
	log.Debug(ctx, "create new rabbitmq server", log.Args{"queue": queue})

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
				userId := msg.Headers["UserId"].(string)
				ctx := qqcontext.WithUserIdValue(context.Background(), userId)

				err = s.handleRawMessage(ctx, msg.Body, msg.CorrelationId, msg.ReplyTo)
				if err != nil {
					log.Error(ctx, "failed to handle message", log.Args{"error": err})
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
	var baseMessage protocol.BaseMessage
	err := json.Unmarshal(body, &baseMessage)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	switch baseMessage.Name {
	case "add":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(addMessage protocol.AddMessage) protocol.AddReplyMessage {
				entity := qqserver.FromAddMessage(addMessage)
				return qqserver.ToAddReplyMessage(s.service.Add(ctx, entity))
			})
		if err != nil {
			return fmt.Errorf("failed to handle add message: %w", err)
		}

	case "remove":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(removeMessage protocol.RemoveMessage) protocol.RemoveReplyMessage {
				key := qqserver.FromRemoveMessage(removeMessage)
				return qqserver.ToRemoveReplyMessage(s.service.Remove(ctx, key))
			})
		if err != nil {
			return fmt.Errorf("failed to handle remove message: %w", err)
		}

	case "get":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(getMessage protocol.GetMessage) protocol.GetReplyMessage {
				key := qqserver.FromGetMessage(getMessage)
				return qqserver.ToGetReplyMessage(s.service.Get(ctx, key))
			})
		if err != nil {
			return fmt.Errorf("failed to handle get message: %w", err)
		}

	case "get all":
		err = handleMessage(ctx, s, body, corrId, replyTo,
			func(getAllMessage protocol.GetAllMessage) protocol.GetAllReplyMessage {
				return qqserver.ToGetAllReplyMessage(s.service.GetAll(ctx))
			})
		if err != nil {
			return fmt.Errorf("failed to handle add message: %w", err)
		}
	}

	return nil
}
