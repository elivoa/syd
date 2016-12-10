package auth

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/service"
)

type Authorize struct {
	core.Page
}

func (p *Authorize) SetupRender() *exit.Exit {
	err := service.Auth.Srv.HandleAuthorizeRequest(p.W, p.R)
	if err != nil {
		return exit.Error(err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return nil
}
