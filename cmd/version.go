package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gobuffalo/buffalo-goth/goth"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of goth",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("goth", goth.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
