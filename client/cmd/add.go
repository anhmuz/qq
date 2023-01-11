package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqcontext"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [flags] <key> <value>",
	Short: "add item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userId, err := rootCmd.Flags().GetString("user_id")
		if err != nil {
			return err
		}

		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			return err
		}

		ctx := qqcontext.WithUserIdValue(cmd.Context(), userId)

		log.Debug(ctx, "add called")

		c, err := rabbitqq.NewClient(ctx, queue)
		if err != nil {
			log.Error(ctx, "failed to create new client", log.Args{"error": err})
			return err
		}

		key := args[0]
		value := args[1]
		added, err := c.Add(ctx, key, value)
		if err != nil {
			log.Error(ctx, "failed to add", log.Args{"error": err, "key": key, "value": value})
			return err
		}

		log.Info(ctx, "add command result", log.Args{"added": added, "key": key, "value": value})

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
