package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqcontext"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [flags] <key>",
	Short: "get item",
	Args:  cobra.ExactArgs(1),
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

		log.Debug(ctx, "get called")

		c, err := rabbitqq.NewClient(ctx, queue)
		if err != nil {
			log.Error(ctx, "failed to create new client", log.Args{"error": err})
			return err
		}

		key := args[0]
		entity, err := c.Get(ctx, key)
		if err != nil {
			log.Error(ctx, "failed to get reply", log.Args{"error": err, "key": key})
			return err
		}

		if entity != nil {
			log.Info(ctx, "get command result", log.Args{"key": key, "entity": *entity})
		} else {
			log.Info(ctx, "entity does not exist", log.Args{"key": key})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
