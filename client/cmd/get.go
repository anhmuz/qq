package cmd

import (
	"fmt"
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
		c := rabbitqq.NewClient(queue)
		key := args[0]
		c.Get(key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
