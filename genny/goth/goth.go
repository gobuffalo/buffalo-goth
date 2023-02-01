package goth

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates
var templates embed.FS

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if len(opts.Providers) == 0 {
		return g, fmt.Errorf("you must specify at least one provider")
	}

	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return g, fmt.Errorf("failed to get subtree of templates: %w", err)
	}
	if err := g.FS(sub); err != nil {
		return g, fmt.Errorf("failed to add subtree: %w", err)
	}

	h := template.FuncMap{
		"downcase": strings.ToLower,
		"upcase":   strings.ToUpper,
	}
	data := map[string]interface{}{
		"providers": opts.Providers,
	}
	t := gogen.TemplateTransformer(data, h)
	g.Transformer(t)

	g.RunFn(func(r *genny.Runner) error {
		path := "actions/app.go"

		f, err := r.FindFile(path)
		if err != nil {
			return fmt.Errorf("setup goth: %w", err)
		}

		f, err = gogen.AddImport(f, "github.com/markbates/goth/gothic")
		if err != nil {
			return fmt.Errorf("could not add import: %w", err)
		}

		expressions := []string{
			"auth := app.Group(\"/auth\")",
			"auth.GET(\"/{provider}\", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))",
			"auth.GET(\"/{provider}/callback\", AuthCallback)",
		}

		f, err = gogen.AddInsideBlock(f, "appOnce.Do(func() {", expressions...)
		if err != nil {
			if strings.Contains(err.Error(), "could not find desired block") {
				// TODO: remove this block some day soon
				// add this block for compatibility with the apps built with
				// the old version of Buffalo CLI (v0.18.8 or older)
				f, err = gogen.AddInsideBlock(f, "if app == nil {", expressions...)
				if err != nil {
					if err != nil {
						return fmt.Errorf("could not add a code block: %w", err)
					} else {
						r.Logger.Warnf("This app was built with CLI v0.18.8 or older. See https://gobuffalo.io/documentation/known-issues/#cli-v0.18.8")
					}
				}
			} else {
				return fmt.Errorf("could not add a code block: %w", err)
			}
		}
		return r.File(f)
	})

	return g, nil
}
