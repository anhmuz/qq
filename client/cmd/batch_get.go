package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqclient"
	"qq/pkg/qqclient/http"
	"qq/pkg/qqclient/rabbitqq"
	"qq/pkg/qqcontext"

	"github.com/spf13/cobra"
)

type keyReply struct {
	key               string
	asyncReplyChannel chan qqclient.AsyncReply[*qqclient.Entity]
}

var batchGetCmd = &cobra.Command{
	Use:   "batch-get [flags] <key1> <key2> ...",
	Short: "get items",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userId, err := rootCmd.Flags().GetString("user_id")
		if err != nil {
			log.Error(cmd.Context(), "failed to get user ID value from command flag ", log.Args{"error": err})
			return err
		}

		ctx := qqcontext.WithUserIdValue(cmd.Context(), userId)

		clientType, err := rootCmd.Flags().GetString("client_type")
		if err != nil {
			log.Error(ctx, "failed to get client type value from command flag ", log.Args{"error": err})
			return err
		}

		var c qqclient.Client

		if clientType == rabbitqq.ClientType {
			queue, err := rootCmd.Flags().GetString("queue")
			if err != nil {
				log.Error(ctx, "failed to get queue value from command flag ", log.Args{"error": err})
				return err
			}

			c, err = rabbitqq.NewClient(ctx, queue)
			if err != nil {
				log.Error(ctx, "failed to create new client", log.Args{"error": err})
				return err
			}
		} else if clientType == http.ClientType {
			c = http.NewClient(ctx)
		}

		log.Debug(ctx, "batch-get called")

		keyReplies := make([]keyReply, 0, len(args))

		for _, key := range args {
			asyncReplyCh, err := c.GetAsync(ctx, key)

			if err != nil {
				continue
			}

			keyReply := keyReply{
				key:               key,
				asyncReplyChannel: asyncReplyCh,
			}

			keyReplies = append(keyReplies, keyReply)
		}

		for _, keyReply := range keyReplies {
			asyncReply := <-keyReply.asyncReplyChannel

			if asyncReply.Err != nil {
				log.Error(ctx, "failed to get reply", log.Args{"error": err, "key": keyReply.key})
				continue
			}

			entity := asyncReply.Result
			if entity == nil {
				log.Info(ctx, "entity does not exist", log.Args{"key": keyReply.key})
				continue
			}

			log.Info(ctx, "batch-get command result", log.Args{"key": keyReply.key, "entity": *entity})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(batchGetCmd)
}
