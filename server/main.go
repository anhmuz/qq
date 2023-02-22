package main

import (
	"context"
	"fmt"
	"qq/pkg/log"
	"qq/repos/cacheqq"
	"qq/repos/qq"
	"qq/server/qqserver/http"
	qqServ "qq/services/qq"
)

const HTTPServerURL = "localhost:8080"

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

	/*server, err := rabbitqqSrv.NewServer(ctx, rabbitqq.RpcQueue, service)
	if err != nil {
		log.Critical(ctx, "failed to create new RabbitMQ server", log.Args{"error": err})
		panic(fmt.Errorf("failed to create new RabbitMQ server: %w", err))
	}*/

	server, err := http.NewServer(ctx, HTTPServerURL, service)
	if err != nil {
		log.Critical(ctx, "failed to create new http server", log.Args{"error": err})
		panic(fmt.Errorf("failed to create new http server: %w", err))
	}

	err = server.Serve()
	if err != nil {
		log.Critical(ctx, "failed to serve", log.Args{"error": err})
		panic(fmt.Errorf("failed to serve: %w", err))
	}
}
