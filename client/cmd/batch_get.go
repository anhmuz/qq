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

		replyErrorChannels := make([]chan rabbitqq.ReplyError, len(args))

		for i, key := range args {
			i := i
			chReplyError, err := c.GetAsync(key)

			if err != nil {
				log.Println(err)
			} else {
				replyErrorChannels[i] = chReplyError
			}
		}

		for _, replyErrorCh := range replyErrorChannels {
			replyError := <-replyErrorCh

			if replyError.Err != nil {
				log.Println(replyError.Err)
			} else {
				log.Printf("value: %+v", *replyError.Reply.(*rabbitqq.GetReplyMessage).Value)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(batchGetCmd)
}
