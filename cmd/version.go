package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// use this command to build with version
// `go build -ldflags "-X cpenv/cmd.Version=v1.2.3" -o cpenv ./main.go`
var Version = "v0.0.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cpenv",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", Version)
	},
}
