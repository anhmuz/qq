package rabbitqq

import (
	"fmt"
)

type Client interface {
	Add(key string, value string)
	Remove(key string)
	Get(key string)
	GetAll()
}

type client struct {
	//r rabbitmq client
}

func NewClient() Client {
	c := client{}
	return c
}

func (c client) Add(key string, value string) {
	fmt.Printf("rabbitmq client: add key:%v, value:%v\n", key, value)
}

func (c client) Remove(key string) {
	fmt.Printf("rabbitmq client: remove key:%v\n", key)

}

func (c client) Get(key string) {
	fmt.Printf("rabbitmq client: get key:%v\n", key)

}

func (c client) GetAll() {
	fmt.Println("rabbitmq client: get all")
}
