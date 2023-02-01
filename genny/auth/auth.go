package auth

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/meta"

	"github.com/gobuffalo/buffalo-goth/genny/goth"
)

//go:embed templates
var templates embed.FS

func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}
	g, err := goth.New(&goth.Options{
		Providers: opts.Providers,
	})
	if err != nil {
		return gg, err
	}
	gg.Add(g)

	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return gg, fmt.Errorf("failed to get subtree of templates: %w", err)
	}

	g = genny.New()

	if err := g.FS(sub); err != nil {
		return gg, fmt.Errorf("failed to add subtree: %w", err)
	}

	h := template.FuncMap{
		"downcase": strings.ToLower,
		"upcase":   strings.ToUpper,
	}
	data := map[string]interface{}{
		"providers": opts.Providers,
		"app":       meta.New("."),
	}
	t := gogen.TemplateTransformer(data, h)
	g.Transformer(t)

	cmd := exec.Command("buffalo", "db", "generate", "model", "user", "name", "email:nulls.String", "provider", "provider_id")
	g.Command(cmd)

	g.RunFn(func(r *genny.Runner) error {
		path := "actions/app.go"

		f, err := r.FindFile(path)
		if err != nil {
			return fmt.Errorf("setup goth: %w", err)
		}

		bb := &bytes.Buffer{}
		for _, line := range strings.Split(f.String(), "\n") {
			if strings.Contains(line, "buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))") {
				expressions := []string{
					"\t\tapp.Use(SetCurrentUser)",
					"\t\tapp.Use(Authorize)",
					"\t\tapp.Middleware.Skip(Authorize, HomeHandler)",
					"\t\tbah := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)",
					"\t\tauth.GET(\"/{provider}\", bah)",
					"\t\tauth.DELETE(\"\", AuthDestroy)",
					"\t\tauth.Middleware.Skip(Authorize, bah, AuthCallback)",
				}
				line = strings.Join(expressions, "\n")
			}
			bb.WriteString(line + "\n")
		}
		f = genny.NewFile(path, bb)
		return r.File(f)
	})

	gg.Add(g)
	return gg, nil
}
