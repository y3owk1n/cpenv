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

type vaultCommand struct{}

func newVaultCmd() *cobra.Command {
	vc := &vaultCommand{}

	return &cobra.Command{
		Use:               "vault",
		Short:             "Open vault in finder",
		Aliases:           []string{"v", "vault"},
		PersistentPreRunE: vc.preRun,
		RunE:              vc.run,
	}
}

func (vc *vaultCommand) preRun(cmd *cobra.Command, args []string) error {
	logrus.WithField("args", args).Debug("Starting vault preRun command")

	// Check if the config file is used and log it.
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init` first"))
		os.Exit(0)
	}
	logrus.Debugf("Using config file: %s", configPath)

	// Retrieve and log the vault directory from the config.
	vaultDir := viper.GetString("vault_dir")
	logrus.Debugf("Vault directory from config: %s", vaultDir)

	// Get the full vault directory and log the result.
	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		logrus.Errorf("Failed to get vault directory: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Resolved full vault directory: %s", vaultDirFull)

	// Set context values and log that the context is updated.
	ctx := cmd.Context()
	ctx = context.WithValue(ctx, ConfigKey, configPath)
	ctx = context.WithValue(ctx, VaultKey, vaultDirFull)
	cmd.SetContext(ctx)
	logrus.Debugf("Context set with ConfigKey=%s and VaultKey=%s", configPath, vaultDirFull)

	return nil
}

func (vc *vaultCommand) run(cmd *cobra.Command, args []string) error {
	logrus.WithField("args", args).Debug("Starting vault run command")

	// Retrieve the vault directory from context and log it.
	vaultDir, ok := cmd.Context().Value(VaultKey).(string)
	if !ok {
		logrus.Error("Vault directory not found in context")
		return fmt.Errorf("vault config not found in context")
	}
	logrus.Debugf("Retrieved vault directory from context: %s", vaultDir)

	// Attempt to open the vault directory in Finder.
	if err := utils.OpenInFinder(vaultDir); err != nil {
		logrus.Errorf("Failed to open vault directory in finder: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Vault directory opened in finder successfully: %s", vaultDir)

	fmt.Println(utils.SuccessMessage.Render("Successfully opened vault in finder."))
	return nil
}

func init() {
	rootCmd.AddCommand(newVaultCmd())
}
