package cmd

import (
	"qq/pkg/log"
	"qq/pkg/qqclient"
	"qq/pkg/qqclient/http"
	"qq/pkg/qqclient/rabbitqq"
	"qq/pkg/qqcontext"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [flags] <key> <value>",
	Short: "add item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userId, err := rootCmd.Flags().GetString("user_id")
		if err != nil {
			log.Error(cmd.Context(), "failed to get user ID value from command flag ", log.Args{"error": err})
			return err
		}

		ctx := qqcontext.WithUserIdValue(cmd.Context(), userId)

		clientType, err := rootCmd.Flags().GetString("client_type")
		if err != nil {
			log.Error(ctx, "failed to get client type value from command flag ", log.Args{"error": err})
			return err
		}

		var c qqclient.Client

		if clientType == rabbitqq.ClientType {
			queue, err := rootCmd.Flags().GetString("queue")
			if err != nil {
				log.Error(ctx, "failed to get queue value from command flag ", log.Args{"error": err})
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

		log.Debug(ctx, "add called")

		key := args[0]
		value := args[1]
		added, err := c.Add(ctx, qqclient.Entity{Key: key, Value: value})
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
