package cmd

import (
	"context"
	"fmt"
	"qq/pkg/log"
	"qq/pkg/qqclient"
	"qq/pkg/qqclient/http"
	"qq/pkg/qqclient/rabbitqq"
	"qq/pkg/qqcontext"
)

func createClient(cmdContext context.Context) (qqclient.Client, context.Context, error) {
	userId, err := rootCmd.Flags().GetString("user_id")
	if err != nil {
		log.Error(cmdContext, "failed to get user ID value from command flag ", log.Args{"error": err})
		return nil, nil, err
	}

	ctx := qqcontext.WithUserIdValue(cmdContext, userId)

	clientType, err := rootCmd.Flags().GetString("client_type")
	if err != nil {
		log.Error(ctx, "failed to get client type value from command flag ", log.Args{"error": err})
		return nil, nil, err
	}

	var client qqclient.Client

	switch clientType {
	case rabbitqq.ClientType:
		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			log.Error(ctx, "failed to get queue value from command flag ", log.Args{"error": err})
			return nil, nil, err
		}

		client, err = rabbitqq.NewClient(ctx, queue)
		if err != nil {
			log.Error(ctx, "failed to create new client", log.Args{"error": err})
			return nil, nil, err
		}

	case http.ClientType:
		client = http.NewClient(ctx)

	default:
		errText := "invalid client type"
		log.Error(ctx, errText)
		return nil, nil, fmt.Errorf(errText)
	}

	return client, ctx, nil
}
