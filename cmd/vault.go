package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

func init() {
	rootCmd.AddCommand(vaultCmd)
}

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Open vault in finder",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig()
		if err != nil {
			utils.Logger.Debug("Failed to load config", "message", err)
			fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
			return
		}

		utils.Logger.Debug("Running vauldCmd", "Vault Directory", config.VaultDir)

		envFilesDirectory, err := core.GetEnvFilesDirectory(config.VaultDir)
		if err != nil {
			utils.Logger.Error("Failed to create env file directory", "message", err)
			return
		}

		utils.Logger.Debugf("Opening vault in finder: %s", envFilesDirectory)

		if err := utils.OpenInFinder(envFilesDirectory); err != nil {
			utils.Logger.Error("Failed to open the directory", "message", err)
		} else {
			fmt.Println(utils.SuccessMessage.Render("Successfully opened the dir in finder."))
		}
	},
}
