package cmd

import (
	"fmt"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [flags] <key>",
	Short: "remove item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("remove called")

		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			return err
		}
		c := rabbitqq.NewClient(queue)
		key := args[0]
		c.Remove(key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
