package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqclient"
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
		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			return err
		}

		userId, err := rootCmd.Flags().GetString("user_id")
		if err != nil {
			return err
		}

		ctx := qqcontext.WithUserIdValue(cmd.Context(), userId)

		log.Debug(ctx, "batch-get called")

		c, err := rabbitqq.NewClient(ctx, queue)
		if err != nil {
			log.Error(ctx, "failed to create new client", log.Args{"error": err})
			return err
		}

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
