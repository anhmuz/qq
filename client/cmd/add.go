package cmd

import (
	"fmt"
	"log"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [flags] <key> <value>",
	Short: "add item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("add called")

		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			return err
		}

		c, err := rabbitqq.NewClient(queue)
		if err != nil {
			return err
		}

		key := args[0]
		value := args[1]
		added, err := c.Add(cmd.Context(), key, value)
		if err != nil {
			return err
		}

		log.Printf("Added: %v", added)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
