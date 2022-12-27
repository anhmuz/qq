package cmd

import (
	"fmt"
	"log"
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

		c, err := rabbitqq.NewClient(queue)
		if err != nil {
			return err
		}

		key := args[0]
		removed, err := c.Remove(cmd.Context(), key)
		if err != nil {
			return err
		}

		log.Printf("Removed: %v", removed)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
