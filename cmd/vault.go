package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

type vaultCommand struct {
	logger *log.Logger
}

func newVaultCmd() *cobra.Command {
	vc := &vaultCommand{
		logger: utils.Logger, // Use the existing logger
	}

	return &cobra.Command{
		Use:   "vault",
		Short: "Open vault in finder",
		RunE:  vc.run,
	}
}

func (vc *vaultCommand) run(cmd *cobra.Command, args []string) error {
	config, err := core.LoadConfig()
	if err != nil {
		vc.logger.Debug("Failed to load config",
			"error", err,
			"suggestion", "Please run `cpenv config init`",
		)
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
		os.Exit(0)
		return err
	}

	vc.logger.Debug("Running vaultCmd",
		"vaultDirectory", config.VaultDir,
	)

	envFilesDirectory, err := core.GetEnvFilesDirectory(config.VaultDir)
	if err != nil {
		vc.logger.Error("Failed to create env file directory",
			"error", err,
		)
		return err
	}

	if err := utils.OpenInFinder(envFilesDirectory); err != nil {
		vc.logger.Error("Failed to open the directory",
			"error", err,
		)
		return err
	}

	fmt.Println(utils.SuccessMessage.Render("Successfully opened the dir in finder."))
	return nil
}
