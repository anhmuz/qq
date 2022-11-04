package rabbitqq

import (
	"fmt"
	"log"
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

var _ Server = &server{}

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

func (s *server) Serve() error {
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
		log.Printf("Received a message: %s", msg.Body)
	}

	return nil
}

func (s server) parse() {

}
