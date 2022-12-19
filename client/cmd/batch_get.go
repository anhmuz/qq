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

		asyncReplyChannels := make([]chan rabbitqq.AsyncReply[rabbitqq.GetReplyMessage], 0, len(args))

		for _, key := range args {
			asyncReplyCh, err := c.GetAsync(key)

			if err != nil {
				log.Println(err)
				continue
			}

			asyncReplyChannels = append(asyncReplyChannels, asyncReplyCh)
		}

		for _, asyncReplyCh := range asyncReplyChannels {
			asyncReply := <-asyncReplyCh

			if asyncReply.Err != nil {
				log.Println(asyncReply.Err)
				continue
			}

			value := asyncReply.Reply.Value
			if value == nil {
				log.Println("value does not exist")
				continue
			}

			log.Printf("value: %+v", *value)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(batchGetCmd)
}
