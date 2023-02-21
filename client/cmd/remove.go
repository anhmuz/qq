package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqclient/rabbitqq"
	"qq/pkg/qqcontext"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [flags] <key>",
	Short: "remove item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userId, err := rootCmd.Flags().GetString("user_id")
		if err != nil {
			log.Error(cmd.Context(), "failed to get user ID value from command flag ", log.Args{"error": err})
			return err
		}

		ctx := qqcontext.WithUserIdValue(cmd.Context(), userId)

		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			log.Error(ctx, "failed to get queue value from command flag ", log.Args{"error": err})
			return err
		}

		log.Debug(ctx, "remove called")

		c, err := rabbitqq.NewClient(ctx, queue)
		if err != nil {
			log.Error(ctx, "failed to create new client", log.Args{"error": err})
			return err
		}

		key := args[0]
		removed, err := c.Remove(ctx, key)
		if err != nil {
			log.Error(ctx, "failed to remove", log.Args{"error": err, "key": key})
			return err
		}

		log.Info(ctx, "remove command result", log.Args{"removed": removed, "key": key})

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
