/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"qq/pkg/rabbitqq"

	"github.com/spf13/cobra"
)

// getAllCmd represents the getAll command
var getAllCmd = &cobra.Command{
	Use:   "getAll",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("get all expects no arguments")
		}

		fmt.Println("getAll called")

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getAllCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
