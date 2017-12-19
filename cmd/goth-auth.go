package cmd

import (
	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo-goth/auth"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

// gothCmd generates a actions/auth.go file configured to the specified providers.
var gothAuthCmd = &cobra.Command{
	Use:   "goth-auth",
	Short: "Generates a full auth implementation use Goth",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify at least one provider")
		}

		g, err := auth.New()
		if err != nil {
			return err
		}
		return g.Run(".", makr.Data{
			"providers":   args,
			"packagePath": envy.CurrentPackage(),
		})
	},
}

func init() {
	RootCmd.AddCommand(gothAuthCmd)
}
