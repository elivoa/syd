package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/service"
)

type OAuth2 struct {
	core.Page
	Title string
}

func (p *OAuth2) SetupRender() *exit.Exit {
	p.R.ParseForm()
	state := p.R.Form.Get("state")
	if state != "xyz" {
		return exit.Error(errors.New("State invalid"))
	}
	code := p.R.Form.Get("code")
	if code == "" {
		return exit.Error(errors.New("Code not found"))
	}
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return exit.Error(err)
	}

	// here is the output. we change this into redirect to homepage.
	// e := json.NewEncoder(p.W)
	// e.SetIndent("", "  ")
	// e.Encode(*token)

	// here we got our token.

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	service.Auth.OAuthUpdateToken(p.W, p.R, token)

	return exit.Redirect("/") // redirect to index.html.
}
