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

type vaultCommand struct {
	logger *log.Logger
}

func newVaultCmd() *cobra.Command {
	vc := &vaultCommand{
		logger: utils.Logger,
	}

	return &cobra.Command{
		Use:               "vault",
		Short:             "Open vault in finder",
		Aliases:           []string{"v", "vault"},
		PersistentPreRunE: vc.preRun,
		RunE:              vc.run,
	}
}

func (vc *vaultCommand) preRun(cmd *cobra.Command, args []string) error {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init` first"))
		os.Exit(0)
	}

	vaultDir := viper.GetString("vault_dir")

	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		vc.logger.Error("Failed to get env file directory",
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

func (vc *vaultCommand) run(cmd *cobra.Command, args []string) error {
	vaultDir, ok := cmd.Context().Value("vault").(string)
	if !ok {
		return fmt.Errorf("config not found in context")
	}

	vc.logger.Debug("Running vaultCmd",
		"vaultDirectory", vaultDir,
	)

	if err := utils.OpenInFinder(vaultDir); err != nil {
		vc.logger.Error("Failed to open the directory",
			"error", err,
		)
		os.Exit(1)
	}

	fmt.Println(utils.SuccessMessage.Render("Successfully opened vault in finder."))
	return nil
}
