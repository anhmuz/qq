package cmd

import (
	"os"
	"path/filepath"
	"qq/pkg/qqclient/rabbitqq"
	"qq/pkg/qqcontext"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "QQ client",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("queue", rabbitqq.RpcQueue, "Queue name")
	rootCmd.PersistentFlags().String("user_id", qqcontext.DefaultUserIdValue, "User ID")
}
