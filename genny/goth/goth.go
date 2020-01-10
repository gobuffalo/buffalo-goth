package goth

import (
	"os/exec"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen/gomods"
	"github.com/gobuffalo/genny/gogen"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if len(opts.Providers) == 0 {
		return g, errors.New("you must specify at least one provider")
	}

	if err := g.Box(packr.NewBox("../goth/templates")); err != nil {
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

	pkg := "github.com/markbates/goth"
	if !gomods.On() {
		pkg += "/..."
	}
	g.Command(exec.Command(genny.GoBin(), "get", pkg))

	return g, nil
}
