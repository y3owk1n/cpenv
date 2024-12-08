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
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type configInitCommand struct {
	logger *log.Logger
}

type configEditCommand struct {
	logger *log.Logger
}

// configCmd represents the init command
var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Config management for cpenv",
	Aliases: []string{"c", "config"},
}

func newConfigInitCommand() *cobra.Command {
	cic := &configInitCommand{
		logger: utils.Logger, // Use the existing logger
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
		logger: utils.Logger, // Use the existing logger
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
	_, err := core.LoadConfig()
	if err != nil {
		cic.logger.Debug("Failed to load config",
			"error", err,
		)

		cic.logger.Debug("Running configInitCmd")

		toInitConfig := &core.Config{VaultDir: ".env-files"}

		err := core.SaveConfig(toInitConfig)
		if err != nil {
			cic.logger.Error("Failed to save config",
				"error", err,
			)
			return err
		}

		fmt.Println(utils.SuccessMessage.Render("Configuration initialized successfully!"))
		return err
	}

	fmt.Println(utils.ErrorMessage.Render("Configuration exists! Use `cpenv config edit` to edit it"))

	return nil
}

func (cec *configEditCommand) preRun(cmd *cobra.Command, args []string) error {
	config, err := core.LoadConfig()
	if err != nil {
		cec.logger.Debug("Failed to load config",
			"error", err,
			"suggestion", "Please run `cpenv config init`",
		)
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
		os.Exit(0)
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), "config", config))

	return nil
}

func (cec *configEditCommand) run(cmd *cobra.Command, args []string) error {
	config, ok := cmd.Context().Value("config").(*core.Config)
	if !ok {
		return fmt.Errorf("config not found in context")
	}

	cec.logger.Debug("Running configEditCommand",
		"vaultDirectory", config.VaultDir,
	)

	if err := utils.OpenInEditor(core.ConfigPath); err != nil {
		cec.logger.Error("Failed to open configuration file in editor",
			"error", err,
		)
		return err
	}

	fmt.Println(utils.SuccessMessage.Render("Successfully opened the configuration file in editor."))

	return nil
}
