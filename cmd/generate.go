package cmd

import (
	"context"

	"github.com/gobuffalo/buffalo-goth/genny/goth"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateOptions = struct {
	*goth.Options
	dryRun bool
}{
	Options: &goth.Options{},
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "goth",
	Short: "generates a actions/auth.go file configured to the specified providers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if generateOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		opts := generateOptions.Options
		opts.Providers = args
		g, err := goth.New(opts)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		g, err = gogen.Fmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&generateOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	gothCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(generateCmd)
}
