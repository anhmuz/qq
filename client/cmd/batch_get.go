package cmd

import (
	"fmt"
	"log"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

var batchGetCmd = &cobra.Command{
	Use:   "batch-get [flags] <key1> <key2> ...",
	Short: "get items",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("batch-get called")

		queue, err := rootCmd.Flags().GetString("queue")
		if err != nil {
			return err
		}

		c, err := rabbitqq.NewClient(queue)
		if err != nil {
			return err
		}

		ch := make(chan error, len(args))
		channels := make([]chan interface{}, len(args))

		for i, key := range args {
			i := i
			go func(key string) {
				chr, err := c.GetAsync(key)
				channels[i] = chr
				ch <- err
			}(key)
		}

		for i := 0; i < len(args); i++ {
			err := <-ch
			if err != nil {
				return err
			}
		}

		for _, chr := range channels {
			log.Printf("reply: %+v", <-chr)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(batchGetCmd)
}
