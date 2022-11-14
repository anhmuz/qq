package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const CallbackQueue = "callback_queue"

type Client interface {
	Add(key string, value string) error
	Remove(key string) error
	Get(key string) (*string, error)
	GetAll() (map[string]string, error)
}

type client struct {
	queue   string
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
}

var _ Client = client{}

func NewClient(queue string) (cl Client, err error) {
	log.Printf("create new rabbitmq client - queue: %v\n ", queue)

	ch, msgs, err := connect(queue)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return client{
		queue:   queue,
		channel: ch,
		msgs:    msgs,
	}, nil
}

func connect(queue string) (*amqp.Channel, <-chan amqp.Delivery, error) {
	conn, err := amqp.Dial(AmqpServerURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
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
		return nil, nil, fmt.Errorf("failed to declare a queue: %w", err)
	}
	_, err = ch.QueueDeclare(
		CallbackQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare a callback queue: %w", err)
	}

	msgs, err := ch.Consume(
		CallbackQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return ch, msgs, nil
}

func (c client) Add(key string, value string) error {
	message := AddMessage{
		Name:  "add",
		Key:   key,
		Value: value,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	err := c.sendMessage(message, corrId)
	if err != nil {
		return fmt.Errorf("failed to send %+v: %w", message, err)
	}

	reply, err := c.receiveMessage(corrId)
	if err != nil {
		return fmt.Errorf("failed to receive reply: %w", err)
	}

	log.Printf("reply: %+v", reply.(AddReplyMessage))

	return nil
}

func (c client) Remove(key string) error {
	message := RemoveMessage{
		Name: "remove",
		Key:  key,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	err := c.sendMessage(message, corrId)
	if err != nil {
		return fmt.Errorf("failed to send %+v: %w", message, err)
	}

	reply, err := c.receiveMessage(corrId)
	if err != nil {
		return fmt.Errorf("failed to receive reply: %w", err)
	}

	log.Printf("reply: %+v", reply.(RemoveReplyMessage))

	return nil
}

func (c client) Get(key string) (*string, error) {
	message := GetMessage{
		Name: "get",
		Key:  key,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	err := c.sendMessage(message, corrId)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	reply, err := c.receiveMessage(corrId)
	if err != nil {
		return nil, fmt.Errorf("failed to receive reply: %w", err)
	}

	log.Printf("reply: %+v", reply.(GetReplyMessage))

	return reply.(GetReplyMessage).Value, nil
}

func (c client) GetAll() (map[string]string, error) {
	message := GetAllMessage{
		Name: "get all",
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	err := c.sendMessage(message, corrId)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	reply, err := c.receiveMessage(corrId)
	if err != nil {
		return nil, fmt.Errorf("failed to receive reply: %w", err)
	}

	log.Printf("reply: %+v", reply.(GetAllReplyMessage))

	return reply.(GetAllReplyMessage).Entities, nil
}

func (c client) sendMessage(message interface{}, corrId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to produce JSON: %w", err)
	}

	err = c.channel.PublishWithContext(ctx,
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       CallbackQueue,
			Body:          jsonMessage,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

func (c client) receiveMessage(corrId string) (interface{}, error) {
	var message interface{}

	for msg := range c.msgs {
		if corrId == msg.CorrelationId {
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				return nil, fmt.Errorf("failed to parse JSON: %w", err)
			}
			break
		}
	}

	return message, nil
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
