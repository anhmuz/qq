package cmd

import (
	"fmt"
	"qq/pkg/log"

	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Use:   "get-all",
	Short: "get all items",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx, err := createClient(cmd.Context())
		if err != nil {
			log.Error(ctx, "failed to create client", log.Args{"error": err})
			return err
		}

		log.Debug(ctx, "get-all called")

		entities, err := client.GetAll(ctx)
		if err != nil {
			log.Error(ctx, "failed to get all", log.Args{"error": err})
			return err
		}

		data := log.Args{}
		for i, entity := range entities {
			key := fmt.Sprintf("entity %v", i+1)
			data[key] = entity
		}

		log.Info(ctx, "get-all command result", data)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAllCmd)
}
