/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

func init() {
	rootCmd.AddCommand(newBackupCommand())
}

type backupCommand struct {
	logger *log.Logger
}

func newBackupCommand() *cobra.Command {
	bc := &backupCommand{
		logger: utils.Logger, // Use the existing logger
	}

	return &cobra.Command{
		Use:   "backup",
		Short: "Backup env file(s) to your vault",
		RunE:  bc.run,
	}
}

func (bc *backupCommand) run(cmd *cobra.Command, args []string) error {
	config, err := core.LoadConfig()
	if err != nil {
		bc.logger.Debug("Failed to load config",
			"error", err,
			"suggestion", "Please run `cpenv config init`",
		)
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
		os.Exit(0)
		return err
	}

	bc.logger.Debug("Running backupCmd",
		"vaultDirectory", config.VaultDir,
	)

	_, err = core.CreateEnvFilesDirectoryIfNotFound(config.VaultDir)
	if err != nil {
		bc.logger.Error("Failed to create env file directory",
			"error", err,
		)
		return err
	}

	err = core.ConfirmCwd()
	if err != nil {
		bc.logger.Error("Failed to confirm cwd",
			"error", err,
		)
		return err
	}

	if err := core.CopyEnvFilesToVault(); err != nil {
		bc.logger.Error("Failed to copy env files to vault",
			"error", err,
		)
		return err
	}

	return nil
}
