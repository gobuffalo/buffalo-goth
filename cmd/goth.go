package cmd

import (
	"errors"

	"github.com/gobuffalo/buffalo-goth/goth"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

// gothCmd generates a actions/auth.go file configured to the specified providers.
var gothCmd = &cobra.Command{
	Use:   "goth",
	Short: "Generates a actions/auth.go file configured to the specified providers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify at least one provider")
		}
		g, err := goth.New()
		if err != nil {
			return err
		}
		return g.Run(".", makr.Data{
			"providers": args,
		})
	},
}

func init() {
	RootCmd.AddCommand(gothCmd)
}
