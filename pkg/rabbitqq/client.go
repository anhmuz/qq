package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const CallbackQueue = "callback_queue"

type Client interface {
	Add(key string, value string) (bool, error)
	Remove(key string) (bool, error)
	Get(key string) (*string, error)
	GetAll() ([]Entity, error)
	GetAsync(key string) (chan interface{}, error)
}

type client struct {
	queue     string
	channel   *amqp.Channel
	msgs      <-chan amqp.Delivery
	mu        sync.Mutex
	waitQueue map[string]PendingReply
}

var _ Client = &client{}

type PendingReply struct {
	ch         chan interface{}
	emptyReply interface{}
}

func NewClient(queue string) (cl Client, err error) {
	log.Printf("create new rabbitmq client - queue: %v\n ", queue)

	ch, msgs, err := connect(queue)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := client{
		queue:     queue,
		channel:   ch,
		msgs:      msgs,
		waitQueue: make(map[string]PendingReply),
	}

	go client.receiveReplies()

	return &client, nil
}

func (c *client) receiveReplies() {
	for msg := range c.msgs {
		c.mu.Lock()
		pendingReply, ok := c.waitQueue[msg.CorrelationId]
		if !ok {
			log.Println("unexpected correlation id")
			continue
		}
		c.mu.Unlock()

		err := json.Unmarshal(msg.Body, pendingReply.emptyReply)
		if err != nil {
			log.Println(fmt.Errorf("failed to parse JSON: %w", err))
		}

		pendingReply.ch <- pendingReply.emptyReply
	}
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

func (c *client) Add(key string, value string) (bool, error) {
	message := AddMessage{
		BaseMessage: BaseMessage{Name: AddMessageName},
		Key:         key,
		Value:       value,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	chr, err := sendMessage(c, message, corrId, &AddReplyMessage{})
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return (<-chr).(*AddReplyMessage).Added, nil
}

func (c *client) Remove(key string) (bool, error) {
	message := RemoveMessage{
		BaseMessage: BaseMessage{Name: RemoveMessageName},
		Key:         key,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	chr, err := sendMessage(c, message, corrId, &RemoveReplyMessage{})
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return (<-chr).(*RemoveReplyMessage).Removed, nil
}

func (c *client) Get(key string) (*string, error) {
	message := GetMessage{
		BaseMessage: BaseMessage{Name: GetMessageName},
		Key:         key,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	chr, err := sendMessage(c, message, corrId, &GetReplyMessage{})
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return (<-chr).(*GetReplyMessage).Value, nil
}

func (c *client) GetAsync(key string) (chan interface{}, error) {
	message := GetMessage{
		BaseMessage: BaseMessage{Name: GetMessageName},
		Key:         key,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	chr, err := sendMessage(c, message, corrId, &GetReplyMessage{})
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return chr, nil
}

func (c *client) GetAll() ([]Entity, error) {
	message := GetAllMessage{
		BaseMessage: BaseMessage{Name: GetAllMessageName},
	}

	log.Printf("rabbitmq client: %+v\n", message)

	corrId := randomString(32)

	chr, err := sendMessage(c, message, corrId, &GetAllReplyMessage{})
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return (<-chr).(*GetAllReplyMessage).Entities, nil
}

func sendMessage[Message any, Reply any](c *client, message Message, corrId string, emptyReply Reply) (chan interface{}, error) {
	c.mu.Lock()
	c.waitQueue[corrId] = PendingReply{
		ch:         make(chan interface{}, 1),
		emptyReply: emptyReply,
	}
	c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to produce JSON: %w", err)
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
		return nil, fmt.Errorf("failed to publish a message: %w", err)
	}

	c.mu.Lock()
	ch := c.waitQueue[corrId].ch
	c.mu.Unlock()

	return ch, nil
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
