package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gobuffalo/buffalo-goth/goth"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of buffalo-goth",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("buffalo-goth", goth.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
