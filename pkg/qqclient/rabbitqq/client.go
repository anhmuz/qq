package rabbitqq

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"qq/pkg/log"
	"qq/pkg/protocol"
	"qq/pkg/qqclient"
	"qq/pkg/qqcontext"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const CallbackQueue = "callback_queue"

type client struct {
	queue         string
	channel       *amqp.Channel
	msgs          <-chan amqp.Delivery
	mu            sync.Mutex
	callbackQueue map[string]callback
}

var _ qqclient.Client = &client{}

type callback func([]byte)

func NewClient(ctx context.Context, queue string) (cl qqclient.Client, err error) {
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

func (c *client) Add(ctx context.Context, entity protocol.Entity) (bool, error) {
	message := protocol.AddMessage{
		BaseMessage: protocol.BaseMessage{Name: protocol.AddMessageName},
		Key:         entity.Key,
		Value:       entity.Value,
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	proc := func(reply protocol.AddReplyMessage) bool {
		return reply.Added
	}

	asyncReplyCh, err := sendMessage(ctx, c, message, proc)
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func (c *client) Remove(ctx context.Context, key string) (bool, error) {
	message := protocol.RemoveMessage{
		BaseMessage: protocol.BaseMessage{Name: protocol.RemoveMessageName},
		Key:         key,
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	proc := func(reply protocol.RemoveReplyMessage) bool {
		return reply.Removed
	}

	asyncReplyCh, err := sendMessage(ctx, c, message, proc)
	if err != nil {
		return false, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func (c *client) Get(ctx context.Context, key string) (*protocol.Entity, error) {
	asyncReplyCh, err := c.GetAsync(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func (c *client) GetAsync(ctx context.Context, key string) (chan qqclient.AsyncReply[*protocol.Entity], error) {
	message := protocol.GetMessage{
		BaseMessage: protocol.BaseMessage{Name: protocol.GetMessageName},
		Key:         key,
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	proc := func(reply protocol.GetReplyMessage) *protocol.Entity {
		if reply.Value == nil {
			return nil
		}

		return &protocol.Entity{
			Key:   key,
			Value: *reply.Value,
		}
	}

	asyncReplyCh, err := sendMessage(ctx, c, message, proc)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	return asyncReplyCh, nil
}

func (c *client) GetAll(ctx context.Context) ([]protocol.Entity, error) {
	message := protocol.GetAllMessage{
		BaseMessage: protocol.BaseMessage{Name: protocol.GetAllMessageName},
	}

	log.Debug(ctx, "rabbitmq client", log.Args{"message": message})

	proc := func(reply protocol.GetAllReplyMessage) []protocol.Entity {
		return reply.Entities
	}

	asyncReplyCh, err := sendMessage(ctx, c, message, proc)
	if err != nil {
		return nil, fmt.Errorf("failed to send %+v: %w", message, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func sendMessage[Message any, Reply any, Result any](
	ctx context.Context,
	c *client,
	message Message,
	proc func(Reply) Result,
) (chan qqclient.AsyncReply[Result], error) {
	ch := make(chan qqclient.AsyncReply[Result], 1)

	callback := func(body []byte) {
		var reply Reply
		err := json.Unmarshal(body, &reply)
		if err != nil {
			err = fmt.Errorf("failed to parse JSON: %w", err)
		}

		asyncReply := qqclient.AsyncReply[Result]{
			Result: proc(reply),
			Err:    err,
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
