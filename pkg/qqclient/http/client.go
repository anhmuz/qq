package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qq/pkg/log"
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

func (c client) Add(ctx context.Context, entity qqclient.Entity) (bool, error) {
	request := PostRequest{
		Entity: entity,
	}

	method := http.MethodPost
	requestURL := fmt.Sprintf("%s/entities", HTTPServerURL)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	responce, statusCode, err := getResponce[PostRequest, PostResponce](ctx, c, &request, method, requestURL)
	if err != nil {
		return false, fmt.Errorf("failed to get responce: %w", err)
	}

	if statusCode != http.StatusCreated {
		return false, fmt.Errorf(responce.Status)
	}

	return responce.Added, nil
}

func (c client) Remove(ctx context.Context, key string) (bool, error) {
	method := http.MethodDelete
	requestURL := fmt.Sprintf("%s/entities/%s", HTTPServerURL, key)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	responce, statusCode, err := getResponce[any, DeleteResponce](ctx, c, nil, method, requestURL)
	if err != nil {
		return false, fmt.Errorf("failed to get responce: %w", err)
	}

	if statusCode != http.StatusOK {
		return false, fmt.Errorf(responce.Status)
	}

	return responce.Removed, nil
}

func (c client) Get(ctx context.Context, key string) (*qqclient.Entity, error) {
	method := http.MethodGet
	requestURL := fmt.Sprintf("%s/entities/%s", HTTPServerURL, key)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	responce, statusCode, err := getResponce[any, GetResponce](ctx, c, nil, method, requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get responce: %w", err)
	}

	if statusCode != http.StatusOK && statusCode != http.StatusNotFound {
		return nil, fmt.Errorf(responce.Status)
	}

	return responce.Entity, nil
}

func (c client) GetAsync(ctx context.Context, key string) (chan qqclient.AsyncReply[*qqclient.Entity], error) {
	ch := make(chan qqclient.AsyncReply[*qqclient.Entity], 1)

	go func() {
		result, err := c.Get(ctx, key)
		if err != nil {
			ch <- qqclient.AsyncReply[*qqclient.Entity]{
				Err: fmt.Errorf("failed to get key %s: %w", key, err),
			}
		} else {
			ch <- qqclient.AsyncReply[*qqclient.Entity]{
				Result: result,
			}
		}
	}()

	return ch, nil
}

func (c client) GetAll(ctx context.Context) ([]qqclient.Entity, error) {
	method := http.MethodGet
	requestURL := fmt.Sprintf("%s/entities", HTTPServerURL)

	log.Debug(ctx, "http client", log.Args{"method": method, "request URL": requestURL})

	responce, statusCode, err := getResponce[any, GetAllResponce](ctx, c, nil, method, requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get responce: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(responce.Status)
	}

	return responce.Entities, nil
}

func getResponce[Request any, Responce any](
	ctx context.Context,
	c client,
	request *Request,
	method string,
	requestURL string,
) (*Responce, int, error) {
	bodyReader := &bytes.Reader{}
	if request != nil {
		jsonRequest, err := json.Marshal(request)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to produce JSON: %w", err)
		}

		bodyReader = bytes.NewReader(jsonRequest)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to produce request: %w", err)
	}

	userId := qqcontext.GetUserIdValue(ctx)
	req.Header.Add("UserId", userId)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, 0, fmt.Errorf("internal server error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read responce body: %w", err)
	}

	var responce Responce
	err = json.Unmarshal(body, &responce)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &responce, resp.StatusCode, nil
}
