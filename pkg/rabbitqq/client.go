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
	GetAsync(key string) (chan AsyncReply, error)
	GetAll() ([]Entity, error)
}

type client struct {
	queue     string
	channel   *amqp.Channel
	msgs      <-chan amqp.Delivery
	mu        sync.Mutex
	waitQueue map[string]waitQueueItem
}

var _ Client = &client{}

type waitQueueItem struct {
	ch         chan AsyncReply
	emptyReply interface{}
}

type AsyncReply struct {
	Reply interface{}
	Err   error
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
		waitQueue: make(map[string]waitQueueItem),
	}

	go client.dispatch()

	return &client, nil
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

	asyncReplyCh, err := sendMessage[AddMessage, AddReplyMessage](c, message)
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.(*AddReplyMessage).Added, asyncReply.Err
}

func (c *client) Remove(key string) (bool, error) {
	message := RemoveMessage{
		BaseMessage: BaseMessage{Name: RemoveMessageName},
		Key:         key,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	asyncReplyCh, err := sendMessage[RemoveMessage, RemoveReplyMessage](c, message)
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.(*RemoveReplyMessage).Removed, asyncReply.Err
}

func (c *client) Get(key string) (*string, error) {
	asyncReplyCh, err := c.GetAsync(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.(*GetReplyMessage).Value, asyncReply.Err
}

func (c *client) GetAsync(key string) (chan AsyncReply, error) {
	message := GetMessage{
		BaseMessage: BaseMessage{Name: GetMessageName},
		Key:         key,
	}

	log.Printf("rabbitmq client: %+v\n", message)

	asyncReplyCh, err := sendMessage[GetMessage, GetReplyMessage](c, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return asyncReplyCh, nil
}

func (c *client) GetAll() ([]Entity, error) {
	message := GetAllMessage{
		BaseMessage: BaseMessage{Name: GetAllMessageName},
	}

	log.Printf("rabbitmq client: %+v\n", message)

	asyncReplyCh, err := sendMessage[GetAllMessage, GetAllReplyMessage](c, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.(*GetAllReplyMessage).Entities, asyncReply.Err
}

func sendMessage[Message any, Reply any](c *client, message Message) (chan AsyncReply, error) {
	corrId := randomString(32)
	var reply Reply
	ch := make(chan AsyncReply, 1)

	c.mu.Lock()
	c.waitQueue[corrId] = waitQueueItem{
		ch:         ch,
		emptyReply: &reply,
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

	return ch, nil
}

func (c *client) dispatch() {
	for msg := range c.msgs {
		c.mu.Lock()
		pendingReply, ok := c.waitQueue[msg.CorrelationId]
		if !ok {
			c.mu.Unlock()
			log.Printf("unexpected correlation id: %s\n", msg.CorrelationId)
			continue
		}
		delete(c.waitQueue, msg.CorrelationId)
		c.mu.Unlock()

		err := json.Unmarshal(msg.Body, pendingReply.emptyReply)
		if err != nil {
			err = fmt.Errorf("failed to parse JSON: %w", err)
		}

		replyError := AsyncReply{
			Reply: pendingReply.emptyReply,
			Err:   err,
		}

		pendingReply.ch <- replyError
	}
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
