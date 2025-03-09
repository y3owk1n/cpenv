/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/y3owk1n/cpenv/core"
)

type contextKey string

var Version = "v0.0.0"

func Execute() error {
	rootCmd := &cobra.Command{
		Version: Version,
		Use:     "cpenv",
		Short:   "A CLI for copy and paste your local .env to right projects faster",
		Example: "cpenv",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	if err := core.InitViper(); err != nil {
		return err
	}

	rootCmd.AddCommand(newVaultCmd())

	rootCmd.AddCommand(newCopyCommand())

	rootCmd.AddCommand(newBackupCommand())

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(newConfigInitCommand())
	configCmd.AddCommand(newConfigEditCommand())

	return rootCmd.ExecuteContext(context.Background())
}
