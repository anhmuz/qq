package cmd

import (
	"qq/pkg/log"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [flags] <key>",
	Short: "remove item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx, err := createClient(cmd.Context())
		if err != nil {
			log.Error(ctx, "failed to create client", log.Args{"error": err})
			return err
		}

		log.Debug(ctx, "remove called")

		key := args[0]
		removed, err := client.Remove(ctx, key)
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
