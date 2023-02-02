package main

import (
	"context"
	"fmt"
	"qq/pkg/log"
	"qq/pkg/rabbitqq"
	"qq/repos/qq"
	"qq/repos/redisqq"
	rabbitqqSrv "qq/server/servers/rabbitqq"
	qqServ "qq/services/qq"
)

func main() {
	ctx := context.Background()

	database, err := qq.NewDatabase()
	if err != nil {
		log.Critical(ctx, "failed to create new qq database", log.Args{"error": err})
		panic(fmt.Errorf("failed to create new qq database: %w", err))
	}

	redisDatabase, err := redisqq.NewDatabase()
	if err != nil {
		log.Critical(ctx, "failed to create new redis database", log.Args{"error": err})
		panic(fmt.Errorf("failed to create new redis database: %w", err))
	}

	service, err := qqServ.NewService(database, redisDatabase)
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
