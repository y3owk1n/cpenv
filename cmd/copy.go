/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type copyCommand struct {
	logger *log.Logger
}

func newCopyCommand() *cobra.Command {
	cc := &copyCommand{
		logger: utils.Logger,
	}

	return &cobra.Command{
		Use:               "copy",
		Short:             "Copy env file(s) to your current project",
		Aliases:           []string{"cp", "copy"},
		PersistentPreRunE: cc.preRun,
		RunE:              cc.run,
	}
}

func (cc *copyCommand) preRun(cmd *cobra.Command, args []string) error {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init` first"))
		os.Exit(0)
	}

	vaultDir := viper.GetString("vault_dir")

	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		cc.logger.Error("Failed to get env file directory",
			"error", err,
		)
		os.Exit(1)
	}

	cmd.SetContext(context.WithValue(cmd.Context(), "config", configPath))
	cmd.SetContext(context.WithValue(cmd.Context(), "vault", vaultDirFull))

	return nil
}

func (cc *copyCommand) run(cmd *cobra.Command, args []string) error {
	vaultDir, ok := cmd.Context().Value("vault").(string)
	if !ok {
		return fmt.Errorf("config not found in context")
	}

	cc.logger.Debug("Running copyCmd",
		"vaultDirectory", vaultDir,
	)

	directories, err := core.GetProjectsList(vaultDir)
	if err != nil {
		cc.logger.Error("Failed to get project lists",
			"error", err,
		)
		os.Exit(1)
	}

	directory, err := core.SelectProject(directories)
	if err != nil {
		cc.logger.Error("Failed to select project",
			"error", err,
		)
		os.Exit(1)
	}

	if err := core.CopyEnvFilesToProject(directory, "", vaultDir); err != nil {
		cc.logger.Error("Failed to copy env files to project",
			"error", err,
		)
		os.Exit(1)
	}

	return nil
}
