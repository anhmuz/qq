package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqclient"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [flags] <key> <value>",
	Short: "add item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx, err := CreateClient(cmd.Context())
		if err != nil {
			log.Error(ctx, "failed to create client", log.Args{"error": err})
			return err
		}

		log.Debug(ctx, "add called")

		key := args[0]
		value := args[1]
		added, err := client.Add(ctx, qqclient.Entity{Key: key, Value: value})
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
