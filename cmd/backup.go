/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type backupCommand struct {
	logger *log.Logger
}

func newBackupCommand() *cobra.Command {
	bc := &backupCommand{
		logger: utils.Logger,
	}

	return &cobra.Command{
		Use:               "backup",
		Short:             "Backup env file(s) to your vault",
		Aliases:           []string{"bk", "backup"},
		PersistentPreRunE: bc.preRun,
		RunE:              bc.run,
	}
}

func (bc *backupCommand) preRun(cmd *cobra.Command, args []string) error {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init` first"))
		os.Exit(0)
	}

	vaultDir := viper.GetString("vault_dir")

	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		bc.logger.Error("Failed to get env file directory",
			"error", err,
		)
		os.Exit(1)
	}

	cmd.SetContext(context.WithValue(cmd.Context(), ConfigKey, configPath))
	cmd.SetContext(context.WithValue(cmd.Context(), VaultKey, vaultDirFull))

	return nil
}

func (bc *backupCommand) run(cmd *cobra.Command, args []string) error {
	vaultDir, ok := cmd.Context().Value(VaultKey).(string)
	if !ok {
		return fmt.Errorf("config not found in context")
	}

	bc.logger.Debug("Running backupCmd",
		"vaultDirectory", vaultDir,
	)

	err := core.ConfirmCwd()
	if err != nil {
		bc.logger.Error("Failed to confirm cwd",
			"error", err,
		)
		os.Exit(1)
	}

	action := func() {
		if err := core.CopyEnvFilesToVault(vaultDir); err != nil {
			bc.logger.Error("Failed to copy env files to vault",
				"error", err,
			)
			os.Exit(1)
		}
	}

	_ = spinner.New().
		Title(fmt.Sprintf("Backing up to %s", vaultDir)).
		Action(action).
		Run()

	return nil
}

func init() {
	rootCmd.AddCommand(newBackupCommand())
}
