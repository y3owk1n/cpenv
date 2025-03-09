/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type configInitCommand struct {
	logger *log.Logger
}

type configEditCommand struct {
	logger *log.Logger
}

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Config management for cpenv",
	Aliases: []string{"c", "config"},
}

func newConfigInitCommand() *cobra.Command {
	cic := &configInitCommand{
		logger: utils.Logger,
	}

	return &cobra.Command{
		Use:     "init",
		Short:   "Initialize a config for cpenv to work",
		Aliases: []string{"i", "init"},
		RunE:    cic.run,
	}
}

func newConfigEditCommand() *cobra.Command {
	cec := &configEditCommand{
		logger: utils.Logger,
	}

	return &cobra.Command{
		Use:               "edit",
		Short:             "Edit the cpenv config with $EDITOR",
		Aliases:           []string{"e", "edit"},
		PersistentPreRunE: cec.preRun,
		RunE:              cec.run,
	}
}

func (cic *configInitCommand) run(cmd *cobra.Command, args []string) error {
	if viper.ConfigFileUsed() != "" {
		fmt.Println(utils.ErrorMessage.Render("Configuration exists! Use `cpenv config edit` to edit it"))
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".config", "cpenv", "cpenv.yaml")

	defaultVaultDir := ".env-files"

	viper.Set("vault_dir", defaultVaultDir)

	if err := viper.WriteConfigAs(configPath); err != nil {
		cic.logger.Error("Failed to save config",
			"error", err,
		)
		return err
	}

	fmt.Println(utils.SuccessMessage.Render("Configuration initialized successfully!"))

	if _, err := core.CreateVaultIfNotFound(defaultVaultDir); err != nil {
		cic.logger.Error("Failed to create vault directory",
			"error", err,
		)
		return err
	}

	return nil
}

func (cec *configEditCommand) preRun(cmd *cobra.Command, args []string) error {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init` first"))
		os.Exit(0)
	}

	vaultDir := viper.GetString("vault_dir")

	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		cec.logger.Error("Failed to get env file directory",
			"error", err,
		)
		os.Exit(1)
	}

	configKey := contextKey("config")
	vaultKey := contextKey("vault")

	cmd.SetContext(context.WithValue(cmd.Context(), configKey, configPath))
	cmd.SetContext(context.WithValue(cmd.Context(), vaultKey, vaultDirFull))

	return nil
}

func (cec *configEditCommand) run(cmd *cobra.Command, args []string) error {
	configPath, ok := cmd.Context().Value("config").(string)
	if !ok {
		return fmt.Errorf("config not found in context")
	}

	cec.logger.Debug("Running configEditCommand",
		"configPath", configPath,
	)

	if err := utils.OpenInEditor(configPath); err != nil {
		cec.logger.Error("Failed to open configuration file in editor",
			"error", err,
		)
		os.Exit(1)
	}

	fmt.Println(utils.SuccessMessage.Render("Successfully opened the configuration file in editor."))

	if err := viper.ReadInConfig(); err != nil {
		cec.logger.Error("Failed to reload config",
			"error", err,
		)
		os.Exit(1)
	}

	vaultDir := viper.GetString("vault_dir")

	if _, err := core.CreateVaultIfNotFound(vaultDir); err != nil {
		cec.logger.Error("Failed to create vault directory",
			"error", err,
		)
		os.Exit(1)
	}

	return nil
}
