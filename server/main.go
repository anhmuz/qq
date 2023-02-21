package main

import (
	"context"
	"fmt"
	"qq/pkg/log"
	"qq/pkg/qqclient/rabbitqq"
	"qq/repos/cacheqq"
	"qq/repos/qq"
	rabbitqqSrv "qq/server/qqserver/rabbitqq"
	qqServ "qq/services/qq"
)

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

	server, err := rabbitqqSrv.NewServer(ctx, rabbitqq.RpcQueue, service)
	if err != nil {
		log.Critical(ctx, "failed to create new RabbitMQ server", log.Args{"error": err})
		panic(fmt.Errorf("failed to create new RabbitMQ server: %w", err))
	}

	err = server.Serve()
	if err != nil {
		log.Critical(ctx, "failed to serve", log.Args{"error": err})
		panic(fmt.Errorf("failed to serve: %w", err))
	}
}
