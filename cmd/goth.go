package cmd

import (
	"github.com/spf13/cobra"
)

// gothCmd represents the goth command
var gothCmd = &cobra.Command{
	Use:   "goth",
	Short: "description about this plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gothCmd)
}
