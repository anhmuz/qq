package rabbitqq

import (
	"context"
	"fmt"
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
	s := fmt.Sprintf("add key: %s, value: %s", key, value)
	err := c.send(s)
	return err
}

func (c client) Remove(key string) error {
	fmt.Printf("rabbitmq client:  key:%v\n", key)
	s := fmt.Sprintf("remove key: %s", key)
	err := c.send(s)
	return err
}

func (c client) Get(key string) (*string, error) {
	fmt.Printf("rabbitmq client: get key:%v\n", key)
	s := fmt.Sprintf("get key: %s", key)
	err := c.send(s)
	return nil, err
}

func (c client) GetAll() (map[string]string, error) {
	fmt.Println("rabbitmq client: get all")
	err := c.send("get all")
	return nil, err
}

func (c client) send(item string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.channel.PublishWithContext(ctx,
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(item),
		})
	if err != nil {
		return fmt.Errorf("Failed to request: %w", err)
	}
	return nil
}
