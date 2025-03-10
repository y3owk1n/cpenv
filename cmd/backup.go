/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type backupCommand struct{}

func newBackupCommand() *cobra.Command {
	bc := &backupCommand{}

	return &cobra.Command{
		Use:               "backup",
		Short:             "Backup env file(s) to your vault",
		Aliases:           []string{"bk", "backup"},
		PersistentPreRunE: bc.preRun,
		RunE:              bc.run,
	}
}

func (bc *backupCommand) preRun(cmd *cobra.Command, args []string) error {
	logrus.WithField("args", args).Debug("Starting backup command preRun")

	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init` first"))
		os.Exit(0)
	}
	logrus.Debugf("Using config file: %s", configPath)

	vaultDir := viper.GetString("vault_dir")
	logrus.Debugf("Vault directory from config: %s", vaultDir)

	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		logrus.Errorf("Failed to get env file directory: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Resolved full vault directory: %s", vaultDirFull)

	// Set context values for later retrieval
	ctx := cmd.Context()
	ctx = context.WithValue(ctx, ConfigKey, configPath)
	ctx = context.WithValue(ctx, VaultKey, vaultDirFull)
	cmd.SetContext(ctx)
	logrus.Debugf("Context set with ConfigKey=%s and VaultKey=%s", configPath, vaultDirFull)

	return nil
}

func (bc *backupCommand) run(cmd *cobra.Command, args []string) error {
	logrus.WithField("args", args).Debug("Starting backup command run")

	vaultDir, ok := cmd.Context().Value(VaultKey).(string)
	if !ok {
		logrus.Error("Vault directory not found in context")
		return fmt.Errorf("config not found in context")
	}
	logrus.Debugf("Retrieved vault directory from context: %s", vaultDir)

	// Confirm that the current working directory is valid
	if err := core.ConfirmCwd(); err != nil {
		logrus.Errorf("Failed to confirm current working directory: %v", err)
		os.Exit(1)
	}
	logrus.Debug("Current working directory confirmed")

	action := func() {
		logrus.Debugf("Starting backup action: copying env files to vault at %s", vaultDir)
		if err := core.CopyEnvFilesToVault(vaultDir); err != nil {
			logrus.Errorf("Failed to copy env files to vault: %v", err)
			os.Exit(1)
		}
		logrus.Debug("Env files successfully backed up to vault")
	}

	spinnerMessage := fmt.Sprintf("Backing up to %s", vaultDir)
	logrus.Debugf("Starting spinner with message: %s", spinnerMessage)
	_ = spinner.New().
		Title(spinnerMessage).
		Action(action).
		Run()
	logrus.Debug("Spinner action completed")

	return nil
}

func init() {
	rootCmd.AddCommand(newBackupCommand())
}
