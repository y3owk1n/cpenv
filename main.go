/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/y3owk1n/cpenv/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
