/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cpenv/core"
	"cpenv/utils"
	"fmt"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy env file(s) to your current project",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := core.LoadConfig()
		if err != nil {
			utils.Logger.Debug("Failed to load config", "message", err)
			fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
			return
		}

		utils.Logger.Debug("Running copyCmd", "Vault Directory", config.VaultDir)

		_, err = core.CreateEnvFilesDirectoryIfNotFound(config.VaultDir)
		if err != nil {
			utils.Logger.Error("Failed to create env file directory", "message", err)
			return
		}

		directories, err := core.GetProjectsList()

		directory, err := core.SelectProject(directories)
		if err != nil {
			utils.Logger.Error("Failed to select project", "message", err)
			return
		}

		err = core.CopyEnvFilesToProject(directory, "")
		if err != nil {
			utils.Logger.Error("Failed to copy env files to project", "message", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
