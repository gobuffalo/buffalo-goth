package goth

import (
	"embed"
	"io/fs"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/pkg/errors"
)

//go:embed templates
var templates embed.FS

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if len(opts.Providers) == 0 {
		return g, errors.New("you must specify at least one provider")
	}

	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return g, errors.WithStack(err)
	}
	if err := g.FS(sub); err != nil {
		return g, errors.WithStack(err)
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
			return errors.WithStack(err)
		}

		f, err = gogen.AddImport(f, "github.com/markbates/goth/gothic")
		if err != nil {
			return errors.WithStack(err)
		}

		expressions := []string{
			"auth := app.Group(\"/auth\")",
			"auth.GET(\"/{provider}\", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))",
			"auth.GET(\"/{provider}/callback\", AuthCallback)",
		}
		f, err = gogen.AddInsideBlock(f, "if app == nil {", expressions...)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	})

	return g, nil
}
