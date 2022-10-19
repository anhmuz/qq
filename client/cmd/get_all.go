package cmd

import (
	"fmt"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Use:   "get-all",
	Short: "get all items",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("get-all called")

		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			return err
		}
		c := rabbitqq.NewClient(queue)
		c.GetAll()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAllCmd)
}
