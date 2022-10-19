package cmd

import (
	"fmt"
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
		c := rabbitqq.NewClient(queue)
		key := args[0]
		value := args[1]
		c.Add(key, value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
