package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"log"
)

type OAuth2 struct {
	core.Page
	Title string
}

func (p *OAuth2) SetupRender() *exit.Exit {
	fmt.Println("sljdflsjdflsjdflkjsldkfjlskdjflksjd")
	p.R.ParseForm()
	log.Println("---------1")
	state := p.R.Form.Get("state")
	if state != "xyz" {
		log.Println("---------2 Error State invalid")
		return exit.Error(errors.New("State invalid"))
	}
	code := p.R.Form.Get("code")
	if code == "" {
		log.Println("---------3 Error Code not found")
		return exit.Error(errors.New("Code not found"))
	}
	log.Println("---------4", context.Background(), code)
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Println("---------5 - err.Error()", token, err.Error())
		return exit.Error(err)
	}
	log.Println("---------6")
	e := json.NewEncoder(p.W)
	e.SetIndent("", "  ")
	e.Encode(*token)
	log.Println("------7---------", e)
	return nil
}
