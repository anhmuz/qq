package cmd

import (
	"qq/pkg/log"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [flags] <key>",
	Short: "get item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx, err := createClient(cmd.Context())
		if err != nil {
			log.Error(ctx, "failed to create client", log.Args{"error": err})
			return err
		}

		log.Debug(ctx, "get called")

		key := args[0]
		entity, err := client.Get(ctx, key)
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
