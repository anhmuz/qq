package cmd

import (
	"fmt"
	"log"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [flags] <key>",
	Short: "get item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("get called")

		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			return err
		}

		c, err := rabbitqq.NewClient(queue)
		if err != nil {
			return err
		}

		key := args[0]
		value, err := c.Get(cmd.Context(), key)
		if err != nil {
			return err
		}

		log.Printf("Value: %v", value)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
