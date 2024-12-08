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

type copyCommand struct {
	logger *log.Logger
}

func newCopyCommand() *cobra.Command {
	cc := &copyCommand{
		logger: utils.Logger, // Use the existing logger
	}

	return &cobra.Command{
		Use:   "copy",
		Short: "Copy env file(s) to your current project",
		RunE:  cc.run,
	}
}

func (cc *copyCommand) run(cmd *cobra.Command, args []string) error {
	config, err := core.LoadConfig()
	if err != nil {
		cc.logger.Debug("Failed to load config",
			"error", err,
			"suggestion", "Please run `cpenv config init`",
		)
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
		os.Exit(0)
		return err
	}

	cc.logger.Debug("Running copyCmd",
		"vaultDirectory", config.VaultDir,
	)

	_, err = core.CreateEnvFilesDirectoryIfNotFound(config.VaultDir)
	if err != nil {
		cc.logger.Error("Failed to create env file directory",
			"error", err,
		)
		return err
	}

	directories, err := core.GetProjectsList()
	if err != nil {
		cc.logger.Error("Failed to get project lists",
			"error", err,
		)
		return err
	}

	directory, err := core.SelectProject(directories)
	if err != nil {
		cc.logger.Error("Failed to select project",
			"error", err,
		)
		return err
	}

	if err := core.CopyEnvFilesToProject(directory, ""); err != nil {
		cc.logger.Error("Failed to copy env files to project",
			"error", err,
		)
		return err
	}

	return nil
}
