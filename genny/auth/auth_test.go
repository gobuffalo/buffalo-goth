package auth

import (
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

/**
 * Source: https://www.programming-books.io/essential/go/normalize-newlines-1d3abcf6f17c4186bb9617fa14074e48
 */
func NormalizeNewlines(d string) string {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = strings.ReplaceAll(d, string([]byte{13, 10}), string([]byte{10}))
	// replace CF \r (mac) with LF \n (unix)
	d = strings.ReplaceAll(d, string([]byte{13}), string([]byte{10}))
	return d
}

func Test_Auth(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	g := genny.New()
	g.File(genny.NewFileS("actions/app.go", (appBefore)))
	run.With(g)

	gg, err := New(&Options{
		Providers: []string{"github", "twitter"},
	})
	r.NoError(err)

	gg.With(run)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal("actions/app.go", f.Name())
	r.Equal(NormalizeNewlines(appAfter), NormalizeNewlines(f.String()))

	r.Len(res.Commands, 2)

	c := res.Commands[1]
	r.Equal("buffalo db generate model user name email:nulls.String provider provider_id", strings.Join(c.Args, " "))
}

const appBefore = `package actions

import (
	"github.com/gobuffalo/buffalo"
)

var app *buffalo.App

func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{})
	}

	return app
}`

const appAfter = `package actions

import (
	"github.com/gobuffalo/buffalo"

	"github.com/markbates/goth/gothic"
)

var app *buffalo.App

func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{})
		auth := app.Group("/auth")
		app.Use(SetCurrentUser)
		app.Use(Authorize)
		app.Middleware.Skip(Authorize, HomeHandler)
		bah := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("/{provider}", bah)
		auth.DELETE("", AuthDestroy)
		auth.Middleware.Skip(Authorize, bah, AuthCallback)
		auth.GET("/{provider}/callback", AuthCallback)
	}

	return app
}
`
