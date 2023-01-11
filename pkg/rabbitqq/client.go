package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"qq/pkg/log"
	"qq/pkg/qqcontext"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const CallbackQueue = "callback_queue"

type Client interface {
	Add(ctx context.Context, key string, value string) (bool, error)
	Remove(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (*string, error)
	GetAsync(ctx context.Context, key string) (chan AsyncReply[GetReplyMessage], error)
	GetAll(ctx context.Context) ([]Entity, error)
}

type client struct {
	queue         string
	channel       *amqp.Channel
	msgs          <-chan amqp.Delivery
	mu            sync.Mutex
	callbackQueue map[string]callback
}

var _ Client = &client{}

type callback func([]byte)

type AsyncReply[Reply any] struct {
	Reply Reply
	Err   error
}

func NewClient(ctx context.Context, queue string) (cl Client, err error) {
	log.Debug(ctx, "create new rabbitmq client", log.Args{"queue": queue})

	ch, msgs, err := connect(queue)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := client{
		queue:         queue,
		channel:       ch,
		msgs:          msgs,
		callbackQueue: make(map[string]callback),
	}

	go client.dispatch(ctx)

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

func (c *client) Add(ctx context.Context, key string, value string) (bool, error) {
	message := AddMessage{
		BaseMessage: BaseMessage{Name: AddMessageName},
		Key:         key,
		Value:       value,
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	asyncReplyCh, err := sendMessage[AddMessage, AddReplyMessage](ctx, c, message)
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.Added, asyncReply.Err
}

func (c *client) Remove(ctx context.Context, key string) (bool, error) {
	message := RemoveMessage{
		BaseMessage: BaseMessage{Name: RemoveMessageName},
		Key:         key,
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	asyncReplyCh, err := sendMessage[RemoveMessage, RemoveReplyMessage](ctx, c, message)
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.Removed, asyncReply.Err
}

func (c *client) Get(ctx context.Context, key string) (*string, error) {
	asyncReplyCh, err := c.GetAsync(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.Value, asyncReply.Err
}

func (c *client) GetAsync(ctx context.Context, key string) (chan AsyncReply[GetReplyMessage], error) {
	message := GetMessage{
		BaseMessage: BaseMessage{Name: GetMessageName},
		Key:         key,
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	asyncReplyCh, err := sendMessage[GetMessage, GetReplyMessage](ctx, c, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return asyncReplyCh, nil
}

func (c *client) GetAll(ctx context.Context) ([]Entity, error) {
	message := GetAllMessage{
		BaseMessage: BaseMessage{Name: GetAllMessageName},
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	asyncReplyCh, err := sendMessage[GetAllMessage, GetAllReplyMessage](ctx, c, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Reply.Entities, asyncReply.Err
}

func sendMessage[Message any, Reply any](ctx context.Context, c *client, message Message) (chan AsyncReply[Reply], error) {
	ch := make(chan AsyncReply[Reply], 1)

	callback := func(body []byte) {
		var reply Reply
		err := json.Unmarshal(body, &reply)
		if err != nil {
			err = fmt.Errorf("failed to parse JSON: %w", err)
		}

		asyncReply := AsyncReply[Reply]{
			Reply: reply,
			Err:   err,
		}

		ch <- asyncReply
	}

	corrId := randomString(32)
	userId := qqcontext.GetUserIdValue(ctx)

	c.mu.Lock()
	c.callbackQueue[corrId] = callback
	c.mu.Unlock()

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
			Headers:       amqp.Table{"UserId": userId},
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

func (c *client) dispatch(ctx context.Context) {
	for msg := range c.msgs {
		c.mu.Lock()
		callback, ok := c.callbackQueue[msg.CorrelationId]
		if !ok {
			c.mu.Unlock()
			log.Warning(ctx, "unexpected correlation id:", log.Args{"correlation_id": msg.CorrelationId})
			continue
		}
		delete(c.callbackQueue, msg.CorrelationId)
		c.mu.Unlock()

		callback(msg.Body)
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
