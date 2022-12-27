package cmd

import (
	"fmt"
	"log"
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

		c, err := rabbitqq.NewClient(queue)
		if err != nil {
			return err
		}

		entities, err := c.GetAll(cmd.Context())
		if err != nil {
			return err
		}

		for _, entity := range entities {
			log.Printf("Entity: %v\n", entity)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getAllCmd)
}
