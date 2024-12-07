/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Version = "v0.0.0"

func init() {
}

var rootCmd = &cobra.Command{
	Use:   "cpenv",
	Short: "A CLI for copy and paste your local .env to right projects faster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	Version: Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
