package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"qq/models/rabbitqq"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client interface {
	Add(key string, value string) error
	Remove(key string) error
	Get(key string) (*string, error)
	GetAll() (map[string]string, error)
}

type client struct {
	queue   string
	channel *amqp.Channel
}

var _ Client = client{}

func NewClient(queue string) (cl Client, err error) {
	fmt.Printf("create new rabbitmq client - queue: %v\n ", queue)

	ch, err := connect(queue)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return client{
		queue:   queue,
		channel: ch,
	}, nil
}

func connect(queue string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(AmqpServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	_ = q
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return ch, nil
}

func (c client) Add(key string, value string) error {
	fmt.Printf("rabbitmq client: add key:%v, value:%v\n", key, value)

	command := rabbitqq.Command{
		Name:  "add",
		Key:   key,
		Value: value,
	}

	err := c.send(&command)
	if err != nil {
		return fmt.Errorf("failed to add key %v, value %v: %w", key, value, err)
	}
	return nil
}

func (c client) Remove(key string) error {
	fmt.Printf("rabbitmq client: remove key:%v\n", key)

	command := rabbitqq.Command{
		Name: "remove",
		Key:  key,
	}

	err := c.send(&command)
	if err != nil {
		return fmt.Errorf("failed to remove key %v: %w", key, err)
	}
	return nil
}

func (c client) Get(key string) (*string, error) {
	fmt.Printf("rabbitmq client: get key:%v\n", key)

	command := rabbitqq.Command{
		Name: "get",
		Key:  key,
	}

	err := c.send(&command)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %v: %w", key, err)
	}
	return nil, nil
}

func (c client) GetAll() (map[string]string, error) {
	fmt.Println("rabbitmq client: get all")

	command := rabbitqq.Command{
		Name: "get-all",
	}

	err := c.send(&command)
	if err != nil {
		return nil, fmt.Errorf("failed to get all: %w", err)
	}
	return nil, nil
}

func (c client) send(command *rabbitqq.Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonCommand, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to produce JSON: %w", err)
	}

	err = c.channel.PublishWithContext(ctx,
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonCommand,
		})
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}
	return nil
}
