package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qq/pkg/log"
	"qq/pkg/protocol"
	"qq/pkg/qqclient"
	"qq/pkg/qqcontext"
)

type client struct {
	client *http.Client
}

var _ qqclient.Client = &client{}

func NewClient(ctx context.Context) qqclient.Client {
	log.Debug(ctx, "create new http client")

	return client{
		client: &http.Client{},
	}
}

func (c client) Add(ctx context.Context, entity protocol.Entity) (bool, error) {
	message := protocol.AddMessage{
		BaseMessage: protocol.BaseMessage{Name: protocol.AddMessageName},
		Key:         entity.Key,
		Value:       entity.Value,
	}

	method := http.MethodPost
	requestURL := fmt.Sprintf("%s/entities", HTTPServerURL)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	proc := func(reply protocol.AddReplyMessage) bool {
		return reply.Added
	}

	asyncReplyCh, err := doRequest(ctx, c, &message, method, requestURL, proc)
	if err != nil {
		return false, fmt.Errorf("failed to do post request: %w", err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func (c client) Remove(ctx context.Context, key string) (bool, error) {
	method := http.MethodDelete
	requestURL := fmt.Sprintf("%s/entities/%s", HTTPServerURL, key)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	proc := func(reply protocol.RemoveReplyMessage) bool {
		return reply.Removed
	}

	asyncReplyCh, err := doRequest[protocol.RemoveMessage](ctx, c, nil, method, requestURL, proc)
	if err != nil {
		return false, fmt.Errorf("failed to do delete request: %w", err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func (c client) Get(ctx context.Context, key string) (*protocol.Entity, error) {
	asyncReplyCh, err := c.GetAsync(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func (c client) GetAsync(ctx context.Context, key string) (chan qqclient.AsyncReply[*protocol.Entity], error) {
	method := http.MethodGet
	requestURL := fmt.Sprintf("%s/entities/%s", HTTPServerURL, key)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	proc := func(reply protocol.GetReplyMessage) *protocol.Entity {
		if reply.Value == nil {
			return nil
		}

		return &protocol.Entity{
			Key:   key,
			Value: *reply.Value,
		}
	}

	asyncReplyCh, err := doRequest[protocol.GetMessage](ctx, c, nil, method, requestURL, proc)
	if err != nil {
		return nil, fmt.Errorf("failed to do get request: %w", err)
	}

	return asyncReplyCh, nil
}

func (c client) GetAll(ctx context.Context) ([]protocol.Entity, error) {
	method := http.MethodGet
	requestURL := fmt.Sprintf("%s/entities", HTTPServerURL)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	proc := func(reply protocol.GetAllReplyMessage) []protocol.Entity {
		return reply.Entities
	}

	asyncReplyCh, err := doRequest[protocol.GetAllMessage](ctx, c, nil, method, requestURL, proc)
	if err != nil {
		return nil, fmt.Errorf("failed to do get all request: %w", err)
	}

	asyncReply := <-asyncReplyCh

	return asyncReply.Result, asyncReply.Err
}

func doRequest[Message any, Reply any, Result any](
	ctx context.Context,
	c client,
	message *Message,
	method string,
	requestURL string,
	proc func(Reply) Result,
) (chan qqclient.AsyncReply[Result], error) {
	ch := make(chan qqclient.AsyncReply[Result], 1)

	bodyReader := &bytes.Reader{}
	if message != nil {
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			return nil, fmt.Errorf("failed to produce JSON: %w", err)
		}

		bodyReader = bytes.NewReader(jsonMessage)
	}

	userId := qqcontext.GetUserIdValue(ctx)
	req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to produce request: %w", err)
	}
	req.Header.Add("UserId", userId)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read responce body: %w", err)
	} else if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("resource not found")
	} else if resp.StatusCode == http.StatusInternalServerError {
		err = fmt.Errorf("internal server error: %s", string(body))
	}

	var reply Reply

	if err != nil {
		asyncReply := qqclient.AsyncReply[Result]{
			Result: proc(reply),
			Err:    err,
		}
		ch <- asyncReply

		return ch, nil
	}

	err = json.Unmarshal(body, &reply)
	if err != nil {
		err = fmt.Errorf("failed to parse JSON: %w", err)
	}

	asyncReply := qqclient.AsyncReply[Result]{
		Result: proc(reply),
		Err:    err,
	}

	ch <- asyncReply

	return ch, nil
}
