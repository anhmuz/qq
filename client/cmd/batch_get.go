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

		replyErrorChannels := make([]chan rabbitqq.ReplyError, 0, len(args))

		for _, key := range args {
			replyErrorCh, err := c.GetAsync(key)

			if err != nil {
				log.Println(err)
			} else {
				replyErrorChannels = append(replyErrorChannels, replyErrorCh)
			}
		}

		for _, replyErrorCh := range replyErrorChannels {
			replyError := <-replyErrorCh

			if replyError.Err != nil {
				log.Println(replyError.Err)
				continue
			}

			value := replyError.Reply.(*rabbitqq.GetReplyMessage).Value
			if value != nil {
				log.Printf("value: %+v", *value)
				continue
			}

			log.Printf("value: %+v", value)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(batchGetCmd)
}
