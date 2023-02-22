package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qq/pkg/log"
	httpClient "qq/pkg/qqclient/http"
	"qq/server/qqserver"
	"qq/services/qq"
	"strings"
)

type server struct {
	server  http.Server
	service qq.Service
}

var _ qqserver.Server = server{}

func NewServer(ctx context.Context, url string, service qq.Service) (qqserver.Server, error) {
	log.Debug(ctx, "create new http server")

	mux := http.NewServeMux()
	httpServer := http.Server{
		Addr:    url,
		Handler: mux,
	}

	server := server{
		server:  httpServer,
		service: service,
	}

	mux.HandleFunc("/entities", server.entities)
	mux.HandleFunc("/entities/", server.entity)

	return server, nil
}

func (s server) entities(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		err := handleGetAllRequest(ctx, w, s.service)
		if err != nil {
			log.Error(ctx, "failed to handle get all request", log.Args{"error": err})
		}
	case http.MethodPost:
		err := handlePostRequest(ctx, w, req, s.service)
		if err != nil {
			log.Error(ctx, "failed to handle post request", log.Args{"error": err})
		}
	}

}

func (s server) entity(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.URL.Path, "/entities/")
	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		err := handleGetRequest(ctx, w, key, s.service)
		if err != nil {
			log.Error(ctx, "failed to handle get request", log.Args{"error": err})
		}
	case http.MethodDelete:
		err := handleDeleteRequest(ctx, w, key, s.service)
		if err != nil {
			log.Error(ctx, "failed to handle delete request", log.Args{"error": err})
		}
	}
}

func (s server) Serve() error {
	err := s.server.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to run http server: %w", err)
		}
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}

func handlePostRequest(ctx context.Context, w http.ResponseWriter, req *http.Request, service qq.Service) error {
	var responce httpClient.PostResponce

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to read request body: %w", err)
	}

	var request httpClient.PostRequest

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responce.Status = "invalid JSON request"

		jsonResponce, err := json.Marshal(responce)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to produce JSON: %w", err)
		}

		_, err = w.Write(jsonResponce)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to write response: %w", err)
		}

		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	entity := FromPostRequest(request)
	responce = ToPostResponce(service.Add(ctx, entity))

	w.WriteHeader(http.StatusCreated)

	return writeJsonResponce(w, responce)
}

func handleGetRequest(ctx context.Context, w http.ResponseWriter, key string, service qq.Service) error {
	var responce httpClient.GetResponce

	entity := service.Get(ctx, key)
	responce = ToGetResponce(entity)

	w.WriteHeader(http.StatusOK)

	if entity == nil {
		responce.Status = "not found"
		w.WriteHeader(http.StatusNotFound)
	}

	return writeJsonResponce(w, responce)
}

func handleDeleteRequest(ctx context.Context, w http.ResponseWriter, key string, service qq.Service) error {
	var responce httpClient.DeleteResponce

	responce = ToDeleteResponce(service.Remove(ctx, key))

	w.WriteHeader(http.StatusOK)

	return writeJsonResponce(w, responce)
}

func handleGetAllRequest(ctx context.Context, w http.ResponseWriter, service qq.Service) error {
	var responce httpClient.GetAllResponce

	responce = ToGetallResponce(service.GetAll(ctx))

	w.WriteHeader(http.StatusOK)

	return writeJsonResponce(w, responce)
}

func writeJsonResponce[Responce any](w http.ResponseWriter, responce Responce) error {
	jsonResponce, err := json.Marshal(responce)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to produce JSON: %w", err)
	}

	_, err = w.Write(jsonResponce)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}
