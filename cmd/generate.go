package cmd

import (
	"context"
	"os/exec"

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

// generateCmd represents the goth command
// TODO: rename it, rename the file
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

		g = genny.New()
		gomodtidy := exec.Command("go", "mod", "tidy")
		g.Command(gomodtidy)
		r.With(g)

		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&generateOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	rootCmd.AddCommand(generateCmd)
}
