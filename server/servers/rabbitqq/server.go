package rabbitqq

import (
	"encoding/json"
	"fmt"
	rabbitqqCommand "qq/models/rabbitqq"
	"qq/pkg/rabbitqq"
	"qq/services/qq"

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
	fmt.Printf("create new RabbitMQ server - queue: %v\n ", queue)

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
		command := rabbitqqCommand.Command{}
		err := json.Unmarshal(msg.Body, &command)
		if err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		fmt.Printf("%+v\n", command)
		s.parse(command)
	}

	return nil
}

func (s server) parse(command rabbitqqCommand.Command) {
	switch command.Name {
	case "add":
		s.service.Add(command.Key, command.Value)
	case "remove":
		s.service.Remove(command.Key)
	case "get":
		s.service.Get(command.Key)
	case "get-all":
		s.service.GetAll()
	}
}
