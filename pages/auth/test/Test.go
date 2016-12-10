package test

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"golang.org/x/oauth2"
)

var config = oauth2.Config{
	ClientID:     "222222",
	ClientSecret: "22222222",
	Scopes:       []string{"all"},

	RedirectURL: "http://localhost:8880/auth/test/oauth2",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "http://localhost:8880/auth/authorize",
		TokenURL: "http://localhost:8880/auth/token",
	},

	// RedirectURL: "http://localhost:8880/auth/test/oauth2",
	// Endpoint: oauth2.Endpoint{
	// 	AuthURL:  "http://localhost:9096/authorize",
	// 	TokenURL: "http://localhost:9096/token",
	// },
}

type Test struct {
	core.Page
	Title string
}

func (p *Test) SetupRender() *exit.Exit {
	u := config.AuthCodeURL("xyz")
	fmt.Println("++++++++++++++++++++++++++++++u: ", u)
	return exit.Redirect(u)
}
