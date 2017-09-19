package auth

import (
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
)

// New actions/auth.go file configured to the specified providers.
func New() (*makr.Generator, error) {
	g := makr.New()
	files, err := generators.Find(filepath.Join("github.com", "gobuffalo", "buffalo-goth", "auth"))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

	g.Add(makr.NewCommand(exec.Command("buffalo", "db", "generate", "model", "user", "name", "email:nulls.String", "provider", "provider_id")))

	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			err := generators.AddInsideAppBlock(
				`app.Use(SetCurrentUser)`,
				`app.Use(Authorize)`,
				`app.Middleware.Skip(Authorize, HomeHandler)`,
				`auth := app.Group("/auth")`,
				`bah := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)`,
				`auth.GET("/{provider}", bah)`,
				`auth.GET("/{provider}/callback", AuthCallback)`,
				`auth.DELETE("", AuthDestroy)`,
				`auth.Middleware.Skip(Authorize, bah, AuthCallback)`,
			)
			if err != nil {
				return err
			}
			return generators.AddImport(filepath.Join("actions", "app.go"), "github.com/markbates/goth/gothic")
		},
	})
	g.Add(makr.NewCommand(makr.GoGet("github.com/markbates/goth/...")))
	return g, nil
}
