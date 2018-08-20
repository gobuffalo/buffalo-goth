package cmd

import (
	"fmt"

	"github.com/gobuffalo/buffalo-goth/goth"
	"github.com/spf13/cobra"
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
	gothCmd.AddCommand(versionCmd)
}
