package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 注册
	rootCmd.AddCommand(wordCmd)
	rootCmd.AddCommand(timeCmd)
}