package rabbitqq

import (
	"fmt"
)

type Client interface {
	Add(key string, value string)
	Remove(key string)
	Get(key string) string
	GetAll() map[string]string
}

type client struct {
	queue string
	//r rabbitmq client
}

var _ Client = client{}

func NewClient(queue string) Client {
	fmt.Printf("create new client - queue: %v\n ", queue)
	return client{queue: queue}
}

func (c client) Add(key string, value string) {
	fmt.Printf("rabbitmq client: add key:%v, value:%v\n", key, value)
}

func (c client) Remove(key string) {
	fmt.Printf("rabbitmq client: remove key:%v\n", key)

}

func (c client) Get(key string) string {
	fmt.Printf("rabbitmq client: get key:%v\n", key)
	return ""
}

func (c client) GetAll() map[string]string {
	fmt.Println("rabbitmq client: get all")
	return nil
}
