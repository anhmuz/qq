package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qq/pkg/log"
	"qq/pkg/protocol"
	"qq/pkg/qqcontext"
	"qq/server/qqserver"
	"qq/services/qq"
	"strings"
)

const HTTPServerURL = "localhost:8080"

type server struct {
	server http.Server
}

var _ qqserver.Server = server{}
var _service qq.Service

func NewServer(ctx context.Context, service qq.Service) (qqserver.Server, error) {
	log.Debug(ctx, "create new http server")

	_service = service

	mux := http.NewServeMux()
	mux.HandleFunc("/entities", entities)
	mux.HandleFunc("/entities/", entity)

	httpServer := http.Server{
		Addr:    HTTPServerURL,
		Handler: mux,
	}

	return server{
		server: httpServer,
	}, nil
}

func entities(w http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("UserId")
	ctx := qqcontext.WithUserIdValue(context.Background(), userId)

	switch req.Method {
	case http.MethodGet:
		err := handleGetAllRequest(ctx, w)
		if err != nil {
			log.Error(ctx, "failed to handle get all request", log.Args{"error": err})
		}
	case http.MethodPost:
		err := handlePostRequest(ctx, w, req)
		if err != nil {
			log.Error(ctx, "failed to handle post request", log.Args{"error": err})
		}
	}

}

func entity(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.URL.Path, "/entities/")

	userId := req.Header.Get("UserId")
	ctx := qqcontext.WithUserIdValue(context.Background(), userId)

	switch req.Method {
	case http.MethodGet:
		err := handleGetRequest(ctx, w, key)
		if err != nil {
			log.Error(ctx, "failed to handle get request", log.Args{"error": err})
		}
	case http.MethodDelete:
		err := handleDeleteRequest(ctx, w, key)
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

func handlePostRequest(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	var message protocol.AddMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		err = fmt.Errorf("failed to parse JSON: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	entity := qqserver.FromAddMessage(message)
	reply := qqserver.ToAddReplyMessage(_service.Add(ctx, entity))

	jsonReply, err := json.Marshal(reply)
	if err != nil {
		err = fmt.Errorf("failed to produce JSON: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	_, err = w.Write(jsonReply)
	if err != nil {
		err = fmt.Errorf("failed to write response: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	return nil
}

func handleGetRequest(ctx context.Context, w http.ResponseWriter, key string) error {
	reply := qqserver.ToGetReplyMessage(_service.Get(ctx, key))

	jsonReply, err := json.Marshal(reply)
	if err != nil {
		err = fmt.Errorf("failed to produce JSON: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	_, err = w.Write(jsonReply)
	if err != nil {
		err = fmt.Errorf("failed to write response: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	return nil
}

func handleDeleteRequest(ctx context.Context, w http.ResponseWriter, key string) error {
	reply := qqserver.ToRemoveReplyMessage(_service.Remove(ctx, key))

	jsonReply, err := json.Marshal(reply)
	if err != nil {
		err = fmt.Errorf("failed to produce JSON: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	_, err = w.Write(jsonReply)
	if err != nil {
		err = fmt.Errorf("failed to write response: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	return nil
}

func handleGetAllRequest(ctx context.Context, w http.ResponseWriter) error {
	reply := qqserver.ToGetAllReplyMessage(_service.GetAll(ctx))

	jsonReply, err := json.Marshal(reply)
	if err != nil {
		err = fmt.Errorf("failed to produce JSON: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	_, err = w.Write(jsonReply)
	if err != nil {
		err = fmt.Errorf("failed to write response: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(err)))
		return err
	}

	return nil
}
