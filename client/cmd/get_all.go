package cmd

import (
	"fmt"
	"qq/pkg/log"
	"qq/pkg/qqcontext"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Use:   "get-all",
	Short: "get all items",
	Args:  cobra.NoArgs,
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

		log.Debug(ctx, "get-all called")

		c, err := rabbitqq.NewClient(ctx, queue)
		if err != nil {
			log.Error(ctx, "failed to create new client", log.Args{"error": err})
			return err
		}

		entities, err := c.GetAll(ctx)
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
