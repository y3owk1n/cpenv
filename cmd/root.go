/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/cpenv/core"
)

type contextKey string

var (
	ConfigKey   = contextKey("config")
	VaultKey    = contextKey("vault")
	verbose     bool
	ctx, cancel = context.WithCancel(context.Background())
)

var Version = "v0.0.0"

func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println()
		logrus.Debug("Interrupt received, canceling operations...")
		cancel()
		os.Exit(1)
	}()

	if err := core.InitViper(); err != nil {
		logrus.Fatalf("Failed to init config: %v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:     "cpenv",
	Short:   "A CLI for copy and paste your local .env to right projects faster",
	Version: Version,
}
