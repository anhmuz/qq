package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"qq/models"
	"qq/pkg/qqclient"
	httpClient "qq/pkg/qqclient/http"
	"qq/services/qq"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleGetRequest(t *testing.T) {
	testCases := []struct {
		name          string
		req           *http.Request
		service       qq.ServiceMock
		exp           *qqclient.Entity
		expStatus     string
		expStatusCode int
	}{
		{
			name: "HappyRun",
			req:  httptest.NewRequest(http.MethodGet, "http://localhost:8080/entities/a", nil),
			service: qq.ServiceMock{
				GetMock: func(ctx context.Context, key string) *models.Entity {
					return &models.Entity{Key: "a", Value: "b"}
				},
			},
			exp:           &qqclient.Entity{Key: "a", Value: "b"},
			expStatus:     "",
			expStatusCode: http.StatusOK,
		},
		{
			name: "NotFound",
			req:  httptest.NewRequest(http.MethodGet, "http://localhost:8080/entities/c", nil),
			service: qq.ServiceMock{
				GetMock: func(ctx context.Context, key string) *models.Entity {
					return nil
				},
			},
			exp:           nil,
			expStatus:     http.StatusText(http.StatusNotFound),
			expStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			server := server{
				service: testCase.service,
			}

			w := httptest.NewRecorder()

			mux := newMux(&server)
			mux.ServeHTTP(w, testCase.req)

			resp := w.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responce httpClient.GetResponce

			err = json.Unmarshal(body, &responce)
			assert.NoError(t, err)

			assert.Equal(t, testCase.exp, responce.Entity)
			assert.Equal(t, testCase.expStatus, responce.Status)
			assert.Equal(t, testCase.expStatusCode, resp.StatusCode)
		})
	}
}

func TestHandlePostRequest(t *testing.T) {
	testSuccessRequestJson, err := json.Marshal(httpClient.PostRequest{
		Entity: qqclient.Entity{},
	})
	require.NoError(t, err)

	successBodyReader := bytes.NewReader(testSuccessRequestJson)

	testNotAddedRequestJson, err := json.Marshal(httpClient.PostRequest{
		Entity: qqclient.Entity{},
	})
	require.NoError(t, err)

	notAddedBodyReader := bytes.NewReader(testNotAddedRequestJson)

	testCases := []struct {
		name          string
		req           *http.Request
		service       qq.ServiceMock
		exp           bool
		expStatus     string
		expStatusCode int
	}{
		{
			name: "HappyRun",
			req:  httptest.NewRequest(http.MethodPost, "http://localhost:8080/entities", successBodyReader),
			service: qq.ServiceMock{
				AddMock: func(ctx context.Context, entity models.Entity) bool {
					return true
				},
			},
			exp:           true,
			expStatus:     "",
			expStatusCode: http.StatusCreated,
		},
		{
			name: "NotAdded",
			req:  httptest.NewRequest(http.MethodPost, "http://localhost:8080/entities", notAddedBodyReader),
			service: qq.ServiceMock{
				AddMock: func(ctx context.Context, entity models.Entity) bool {
					return false
				},
			},
			exp:           false,
			expStatus:     "",
			expStatusCode: http.StatusOK,
		},
		{
			name:          "BadRequest",
			req:           httptest.NewRequest(http.MethodPost, "http://localhost:8080/entities", nil),
			service:       qq.ServiceMock{},
			exp:           false,
			expStatus:     http.StatusText(http.StatusBadRequest),
			expStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			server := server{
				service: testCase.service,
			}

			w := httptest.NewRecorder()

			mux := newMux(&server)
			mux.ServeHTTP(w, testCase.req)

			resp := w.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responce httpClient.PostResponce

			err = json.Unmarshal(body, &responce)
			assert.NoError(t, err)

			assert.Equal(t, testCase.exp, responce.Added)
			assert.Equal(t, testCase.expStatus, responce.Status)
			assert.Equal(t, testCase.expStatusCode, resp.StatusCode)
		})
	}
}

func TestHandleDeleteRequest(t *testing.T) {
	testCases := []struct {
		name          string
		req           *http.Request
		service       qq.ServiceMock
		exp           bool
		expStatus     string
		expStatusCode int
	}{
		{
			name: "HappyRun",
			req:  httptest.NewRequest(http.MethodDelete, "http://localhost:8080/entities/a", nil),
			service: qq.ServiceMock{
				RemoveMock: func(ctx context.Context, key string) bool {
					return true
				},
			},
			exp:           true,
			expStatus:     "",
			expStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			server := server{
				service: testCase.service,
			}

			w := httptest.NewRecorder()

			mux := newMux(&server)
			mux.ServeHTTP(w, testCase.req)

			resp := w.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responce httpClient.DeleteResponce

			err = json.Unmarshal(body, &responce)
			assert.NoError(t, err)

			assert.Equal(t, testCase.exp, responce.Removed)
			assert.Equal(t, testCase.expStatus, responce.Status)
			assert.Equal(t, testCase.expStatusCode, resp.StatusCode)
		})
	}
}

func TestHandleGetAllRequest(t *testing.T) {
	testCases := []struct {
		name          string
		req           *http.Request
		service       qq.ServiceMock
		exp           []qqclient.Entity
		expStatus     string
		expStatusCode int
	}{
		{
			name: "HappyRun",
			req:  httptest.NewRequest(http.MethodGet, "http://localhost:8080/entities", nil),
			service: qq.ServiceMock{
				GetAllMock: func(ctx context.Context) []models.Entity {
					return []models.Entity{
						{Key: "a", Value: "b"},
						{Key: "c", Value: "d"},
					}
				},
			},
			exp: []qqclient.Entity{
				{Key: "a", Value: "b"},
				{Key: "c", Value: "d"},
			},
			expStatus:     "",
			expStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			server := server{
				service: testCase.service,
			}

			w := httptest.NewRecorder()

			mux := newMux(&server)
			mux.ServeHTTP(w, testCase.req)

			resp := w.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responce httpClient.GetAllResponce

			err = json.Unmarshal(body, &responce)
			assert.NoError(t, err)

			assert.Equal(t, testCase.exp, responce.Entities)
			assert.Equal(t, testCase.expStatus, responce.Status)
			assert.Equal(t, testCase.expStatusCode, resp.StatusCode)
		})
	}
}
