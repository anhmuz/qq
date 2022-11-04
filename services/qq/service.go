package qq

import (
	"fmt"
)

type Service interface {
	Add(key string, value string)
	Remove(key string)
	Get(key string) *string
	GetAll() map[string]string
}

type service struct {
}

var _ Service = service{}

func NewService() Service {
	fmt.Println("create new service")
	return service{}
}

func (c service) Add(key string, value string) {
	fmt.Printf("service: add key:%v, value:%v\n", key, value)
}

func (c service) Remove(key string) {
	fmt.Printf("service: remove key:%v\n", key)

}

func (c service) Get(key string) *string {
	fmt.Printf("service: get key:%v\n", key)
	return nil
}

func (c service) GetAll() map[string]string {
	fmt.Println("service: get all")
	return nil
}
