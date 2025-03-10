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
		Use:              "vault",
		Short:            "Open vault in finder",
		Aliases:          []string{"v", "vault"},
		PersistentPreRun: vc.preRun,
		Run:              vc.run,
	}
}

func (vc *vaultCommand) preRun(cmd *cobra.Command, args []string) {
	logrus.WithField("args", args).Debug("Starting vault preRun command")

	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		fmt.Printf("%s %s\n", utils.ErrorIcon(), utils.WhiteText("Please run `cpenv config init` first"))
		os.Exit(0)
	}
	logrus.Debugf("Using config file: %s", configPath)

	vaultDir := viper.GetString("vault_dir")
	logrus.Debugf("Vault directory from config: %s", vaultDir)

	vaultDirFull, err := core.GetFullVaultDir(vaultDir)
	if err != nil {
		logrus.Errorf("Failed to get vault directory: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Resolved full vault directory: %s", vaultDirFull)

	ctx := cmd.Context()
	ctx = context.WithValue(ctx, ConfigKey, configPath)
	ctx = context.WithValue(ctx, VaultKey, vaultDirFull)
	cmd.SetContext(ctx)
	logrus.Debugf("Context set with ConfigKey=%s and VaultKey=%s", configPath, vaultDirFull)
}

func (vc *vaultCommand) run(cmd *cobra.Command, args []string) {
	logrus.WithField("args", args).Debug("Starting vault run command")

	vaultDir, ok := cmd.Context().Value(VaultKey).(string)
	if !ok {
		logrus.Error("Vault directory not found in context")
		fmt.Printf("%s %s\n", utils.ErrorIcon(), utils.WhiteText("vault config not found in context"))
		return
	}
	logrus.Debugf("Retrieved vault directory from context: %s", vaultDir)

	if err := utils.OpenInFinder(vaultDir); err != nil {
		logrus.Errorf("Failed to open vault directory in finder: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Vault directory opened in finder successfully: %s", vaultDir)

	fmt.Printf("%s %s\n", utils.SuccessIcon(), utils.WhiteText("Successfully opened vault in finder."))
}

func init() {
	rootCmd.AddCommand(newVaultCmd())
}
