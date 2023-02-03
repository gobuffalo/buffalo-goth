package cmd

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/buffalo-goth/genny/auth"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/spf13/cobra"
)

var authOptions = struct {
	*auth.Options
	dryRun bool
}{
	Options: &auth.Options{},
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "goth-auth",
	Short: "Generates a full auth implementation use Goth",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if authOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		opts := authOptions.Options
		opts.Providers = args
		gg, err := auth.New(opts)
		if err != nil {
			return err
		}
		gg.With(r)

		g, err := gogen.Fmt(r.Root)
		if err != nil {
			return fmt.Errorf("formatting error: %w", err)
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
	authCmd.Flags().BoolVarP(&authOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	rootCmd.AddCommand(authCmd)
}
