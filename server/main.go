package main

import (
	"fmt"
	"qq/pkg/rabbitqq"
	rabbitqqSrv "qq/server/servers/rabbitqq"
	"qq/services/qq"
)

func main() {
	service := qq.NewService()

	server, err := rabbitqqSrv.NewServer(rabbitqq.RpcQueue, service)
	if err != nil {
		panic(fmt.Errorf("failed to create new RabbitMQ server: %w", err))
	}

	err = server.Serve()
	if err != nil {
		panic(fmt.Errorf("failed to serve: %w", err))
	}
}
