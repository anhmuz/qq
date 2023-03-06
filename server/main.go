package main

import (
	"context"
	"fmt"
	"os"
	"qq/pkg/log"
	"qq/pkg/qqclient/rabbitqq"
	"qq/repos/cacheqq"
	"qq/repos/qq"
	"qq/server/qqserver"
	"qq/server/qqserver/http"
	rabbitqqSrv "qq/server/qqserver/rabbitqq"
	qqServ "qq/services/qq"
)

const HTTPServerURL = "localhost:8080"
const HTTPServerType = "http"
const RabbitMQServerType = "rabbitmq"

func main() {
	ctx := context.Background()

	database, err := qq.NewDatabase()
	if err != nil {
		log.Critical(ctx, "failed to create new qq database", log.Args{"error": err})
		panic(fmt.Errorf("failed to create new qq database: %w", err))
	}

	cache := cacheqq.NewRedisCache()

	service, err := qqServ.NewService(database, cache)
	if err != nil {
		log.Critical(ctx, "failed to create new qq service", log.Args{"error": err})
		panic(fmt.Errorf("failed to create new qq service: %w", err))
	}

	serverType := os.Args[1]
	var server qqserver.Server

	switch serverType {
	case HTTPServerType:
		server, err = http.NewServer(ctx, HTTPServerURL, service)
		if err != nil {
			log.Critical(ctx, "failed to create new http server", log.Args{"error": err})
			panic(fmt.Errorf("failed to create new http server: %w", err))
		}
	case RabbitMQServerType:
		server, err = rabbitqqSrv.NewServer(ctx, rabbitqq.RpcQueue, service)
		if err != nil {
			log.Critical(ctx, "failed to create new RabbitMQ server", log.Args{"error": err})
			panic(fmt.Errorf("failed to create new RabbitMQ server: %w", err))
		}
	}

	err = server.Serve()
	if err != nil {
		log.Critical(ctx, "failed to serve", log.Args{"error": err})
		panic(fmt.Errorf("failed to serve: %w", err))
	}
}
