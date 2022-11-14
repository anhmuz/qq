package main

import (
	"fmt"
	"qq/pkg/rabbitqq"
	"qq/repos/qq"
	rabbitqqSrv "qq/server/servers/rabbitqq"
	qqServ "qq/services/qq"
)

func main() {
	database, err := qq.NewDatabase()
	if err != nil {
		panic(fmt.Errorf("failed to create new qq database: %w", err))
	}

	service, err := qqServ.NewService(database)
	if err != nil {
		panic(fmt.Errorf("failed to create new qq service: %w", err))
	}

	server, err := rabbitqqSrv.NewServer(rabbitqq.RpcQueue, service)
	if err != nil {
		panic(fmt.Errorf("failed to create new RabbitMQ server: %w", err))
	}

	err = server.Serve()
	if err != nil {
		panic(fmt.Errorf("failed to serve: %w", err))
	}
}
