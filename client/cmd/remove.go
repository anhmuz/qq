package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqclient"
	"qq/pkg/qqclient/http"
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
			return err
		}
		ctx := qqcontext.WithUserIdValue(cmd.Context(), userId)

		clientType, err := rootCmd.Flags().GetString("client_type")
		if err != nil {
			return err
		}

		var c qqclient.Client
		if clientType == rabbitqq.ClientType {
			queue, err := rootCmd.Flags().GetString("queue")
			if err != nil {
				return err
			}

			c, err = rabbitqq.NewClient(ctx, queue)
			if err != nil {
				log.Error(ctx, "failed to create new client", log.Args{"error": err})
				return err
			}
		} else if clientType == http.ClientType {
			c = http.NewClient(ctx)
		}

		log.Debug(ctx, "remove called")

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
