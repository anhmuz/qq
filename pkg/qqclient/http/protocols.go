package http

import (
	"qq/pkg/qqclient"
)

type PostRequest struct {
	Entity qqclient.Entity `json:"entity"`
}

type PostResponce struct {
	Added  bool   `json:"added"`
	Status string `json:"status"`
}

type DeleteResponce struct {
	Removed bool   `json:"removed"`
	Status  string `json:"status"`
}

type GetResponce struct {
	Entity *qqclient.Entity `json:"entity"`
	Status string           `json:"status"`
}

type GetAllResponce struct {
	Entities []qqclient.Entity `json:"entities"`
	Status   string            `json:"status"`
}
