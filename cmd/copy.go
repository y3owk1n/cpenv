/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type copyCommand struct{}

func newCopyCommand() *cobra.Command {
	cc := &copyCommand{}

	return &cobra.Command{
		Use:               "copy",
		Short:             "Copy env file(s) to your current project",
		Aliases:           []string{"cp", "copy"},
		PersistentPreRunE: cc.preRun,
		RunE:              cc.run,
	}
}

func (cc *copyCommand) preRun(cmd *cobra.Command, args []string) error {
	logrus.WithField("args", args).Debug("Starting copy command preRun")

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

func (cc *copyCommand) run(cmd *cobra.Command, args []string) error {
	logrus.WithField("args", args).Debug("Starting copy command run")

	vaultDir, ok := cmd.Context().Value(VaultKey).(string)
	if !ok {
		logrus.Error("Vault directory not found in context")
		return fmt.Errorf("config not found in context")
	}
	logrus.Debugf("Retrieved vault directory from context: %s", vaultDir)

	directories, err := core.GetProjectsList(vaultDir)
	if err != nil {
		logrus.Errorf("Failed to get project lists: %v", err)
		os.Exit(1)
	}
	logrus.WithField("directories", directories).Debug("Retrieved project list")

	directory, err := core.SelectProject(directories)
	if err != nil {
		logrus.Errorf("Failed to select project: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Selected project directory: %s", directory)

	if err := core.CopyEnvFilesToProject(directory, "", vaultDir); err != nil {
		logrus.Errorf("Failed to copy env files to project: %v", err)
		os.Exit(1)
	}
	logrus.WithFields(logrus.Fields{
		"directory": directory,
		"vaultDir":  vaultDir,
	}).Debug("Successfully copied env files to project")

	return nil
}

func init() {
	rootCmd.AddCommand(newCopyCommand())
}
