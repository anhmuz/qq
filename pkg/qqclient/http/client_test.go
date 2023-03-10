package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qq/pkg/qqclient"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestAdd(t *testing.T) {
	ctx := context.Background()

	testSuccessResponseJson, err := json.Marshal(PostResponce{
		Added: true,
	})
	require.NoError(t, err)

	unmarshal := func(req *http.Request) (qqclient.Entity, error) {
		defer req.Body.Close()

		var request PostRequest
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return request.Entity, err
		}

		err = json.Unmarshal(body, &request)
		return request.Entity, err
	}

	testCases := []struct {
		name       string
		entity     qqclient.Entity
		httpClient *http.Client
		exp        bool
		expErr     error
	}{
		{
			name:   "HappyRun",
			entity: qqclient.Entity{Key: "a", Value: "b"},
			httpClient: NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://localhost:8080/entities", req.URL.String())
				assert.Equal(t, http.MethodPost, req.Method)

				entity, err := unmarshal(req)
				assert.NoError(t, err)

				assert.Equal(t, qqclient.Entity{Key: "a", Value: "b"}, entity)

				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewReader(testSuccessResponseJson)),
				}
			}),
			exp:    true,
			expErr: nil,
		},
		{
			name:   "ServerError",
			entity: qqclient.Entity{Key: "a", Value: "b"},
			httpClient: NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://localhost:8080/entities", req.URL.String())
				assert.Equal(t, http.MethodPost, req.Method)

				entity, err := unmarshal(req)
				assert.NoError(t, err)

				assert.Equal(t, qqclient.Entity{Key: "a", Value: "b"}, entity)

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}
			}),
			exp:    false,
			expErr: fmt.Errorf("failed to get responce: %w", fmt.Errorf("internal server error")),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			client := client{
				client: testCase.httpClient,
			}

			added, err := client.Add(ctx, testCase.entity)
			assert.Equal(t, testCase.exp, added)
			assert.Equal(t, testCase.expErr, err)
		})
	}
}

func TestRemove(t *testing.T) {
	ctx := context.Background()

	testSuccessResponseJson, err := json.Marshal(DeleteResponce{
		Removed: true,
	})
	require.NoError(t, err)

	testCases := []struct {
		name       string
		key        string
		httpClient *http.Client
		exp        bool
		expErr     error
	}{
		{
			name: "HappyRun",
			key:  "a",
			httpClient: NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://localhost:8080/entities/a", req.URL.String())
				assert.Equal(t, http.MethodDelete, req.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(testSuccessResponseJson)),
				}
			}),
			exp:    true,
			expErr: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			client := client{
				client: testCase.httpClient,
			}

			removed, err := client.Remove(ctx, testCase.key)
			assert.Equal(t, testCase.exp, removed)
			assert.Equal(t, testCase.expErr, err)
		})
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	testSuccessResponseJson, err := json.Marshal(GetResponce{
		Entity: &qqclient.Entity{Key: "a", Value: "b"},
	})
	require.NoError(t, err)

	testNotFoundResponseJson, err := json.Marshal(GetResponce{
		//Status: http.StatusText(http.StatusNotFound),
	})
	require.NoError(t, err)

	testCases := []struct {
		name       string
		key        string
		httpClient *http.Client
		exp        *qqclient.Entity
		expErr     error
	}{
		{
			name: "HappyRun",
			key:  "a",
			httpClient: NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://localhost:8080/entities/a", req.URL.String())
				assert.Equal(t, http.MethodGet, req.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(testSuccessResponseJson)),
				}
			}),
			exp:    &qqclient.Entity{Key: "a", Value: "b"},
			expErr: nil,
		},
		{
			name: "NotFound",
			key:  "a",
			httpClient: NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://localhost:8080/entities/a", req.URL.String())
				assert.Equal(t, http.MethodGet, req.Method)

				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(bytes.NewReader(testNotFoundResponseJson)),
				}
			}),
			exp:    nil,
			expErr: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			client := client{
				client: testCase.httpClient,
			}

			entity, err := client.Get(ctx, testCase.key)
			assert.Equal(t, testCase.exp, entity)
			assert.Equal(t, testCase.expErr, err)
		})
	}
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()

	entity1 := qqclient.Entity{Key: "a", Value: "b"}
	entity2 := qqclient.Entity{Key: "c", Value: "d"}

	testSuccessResponseJson, err := json.Marshal(GetAllResponce{
		Entities: []qqclient.Entity{entity1, entity2},
	})
	require.NoError(t, err)

	testCases := []struct {
		name       string
		httpClient *http.Client
		exp        []qqclient.Entity
		expErr     error
	}{
		{
			name: "HappyRun",
			httpClient: NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://localhost:8080/entities", req.URL.String())
				assert.Equal(t, http.MethodGet, req.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(testSuccessResponseJson)),
				}
			}),
			exp:    []qqclient.Entity{entity1, entity2},
			expErr: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			client := client{
				client: testCase.httpClient,
			}

			entities, err := client.GetAll(ctx)
			assert.Equal(t, testCase.exp, entities)
			assert.Equal(t, testCase.expErr, err)
		})
	}
}

func TestGetAsync(t *testing.T) {
	ctx := context.Background()

	testSuccessResponseJson, err := json.Marshal(GetResponce{
		Entity: &qqclient.Entity{Key: "a", Value: "b"},
	})
	require.NoError(t, err)

	testCases := []struct {
		name       string
		key        string
		httpClient *http.Client
		exp        *qqclient.Entity
		expErr     error
	}{
		{
			name: "HappyRun",
			key:  "a",
			httpClient: NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://localhost:8080/entities/a", req.URL.String())
				assert.Equal(t, http.MethodGet, req.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(testSuccessResponseJson)),
				}
			}),
			exp:    &qqclient.Entity{Key: "a", Value: "b"},
			expErr: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			client := client{
				client: testCase.httpClient,
			}

			asyncReplyCh, err := client.GetAsync(ctx, testCase.key)
			asyncReply := <-asyncReplyCh
			assert.Equal(t, *testCase.exp, *asyncReply.Result)
			assert.Equal(t, testCase.expErr, err)
		})
	}
}
