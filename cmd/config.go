/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/y3owk1n/cpenv/core"
	"github.com/y3owk1n/cpenv/utils"
)

// configCmd represents the init command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config management for cpenv",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a config for cpenv to work",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := core.LoadConfig()
		if err != nil {
			toInitConfig := &core.Config{VaultDir: ".env-files"}

			err := core.SaveConfig(toInitConfig)
			if err != nil {
				utils.Logger.Error("Failed to save config", "message", err)
				os.Exit(1)
			}

			fmt.Println(utils.SuccessMessage.Render("Configuration initialized successfully!"))
		} else {
			fmt.Println(utils.ErrorMessage.Render("Configuration exists! Use `cpenv config edit` to edit it"))
		}
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the cpenv config with $EDITOR",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := core.LoadConfig()
		if err != nil {
			utils.Logger.Debug("Failed to load config", "message", err)
			fmt.Println(utils.ErrorMessage.Render("Please run `cpenv config init`"))
			return
		}

		if err := utils.OpenInEditor(core.ConfigPath); err != nil {
			utils.Logger.Error("Failed to open the file in editor", "message", err)
		} else {
			fmt.Println(utils.SuccessMessage.Render("Successfully opened the file in editor."))
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
