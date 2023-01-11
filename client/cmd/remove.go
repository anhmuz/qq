package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqcontext"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [flags] <key>",
	Short: "remove item",
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

		log.Debug(ctx, "remove called", log.Args{})

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
