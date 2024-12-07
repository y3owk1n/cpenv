/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup current project env(s) to vault",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig()
		if err != nil {
			utils.Logger.Debug("Failed to load config", "message", err)
			fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
			return
		}

		utils.Logger.Debug("Running backupCmd", "Vault Directory", config.VaultDir)

		_, err = core.CreateEnvFilesDirectoryIfNotFound(config.VaultDir)
		if err != nil {
			utils.Logger.Error("Failed to create env file directory", "message", err)
			return
		}

		err = core.ConfirmCwd()
		if err != nil {
			utils.Logger.Error("Failed to confirm cwd", "message", err)
			return
		}

		err = core.CopyEnvFilesToVault()
		if err != nil {
			utils.Logger.Error("Failed to copy env files to vault", "message", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
