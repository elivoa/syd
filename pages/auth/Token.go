package auth

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/service"
)

type Token struct {
	core.Page
}

func (p *Token) Activate() {}

func (p *Token) SetupRender() *exit.Exit {
	return nil
}

func (p *Token) OnSuccess() *exit.Exit {
	fmt.Println("55555555555555555555555555555555555555\n success")
	err := service.Auth.Srv.HandleTokenRequest(p.W, p.R)
	if err != nil {
		fmt.Println("6666666666666666666666666666666666666666666666666666666666\n error:\n", err)
		return exit.Error(err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println("6666666666666666666666666666666666666666666666666666666666\n exit")
	return nil
}
