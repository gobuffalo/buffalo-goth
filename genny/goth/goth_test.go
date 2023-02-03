package goth

import (
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/stretchr/testify/require"
)

func Test_Goth(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	g := genny.New()
	g.File(genny.NewFile("actions/app.go", strings.NewReader(appBefore)))
	run.With(g)

	g, err := New(&Options{
		Providers: []string{"github", "twitter"},
	})
	r.NoError(err)
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal("actions/app.go", f.Name())
	r.Equal(appAfter, f.String())

	f = res.Files[1]
	r.Equal("actions/auth.go", f.Name())

	exp := strings.TrimSpace(authAfter)
	act := strings.TrimSpace(f.String())
	act = strings.ReplaceAll(act, "\r\n", "\n")
	r.Equal(exp, act)

	r.Len(res.Commands, 0)

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
		
		// NOTE: this block should go before any resources
		// that need to be protected by buffalo-goth!
		auth := app.Group("/auth")
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		auth.GET("/{provider}/callback", AuthCallback)
		
	}

	return app
}`

const authAfter = `package actions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/twitter"
	)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
	github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
	twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/twitter/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}
	// Do something with the user, maybe register them/sign them in
	return c.Render(http.StatusOK, r.JSON(user))
}
`
