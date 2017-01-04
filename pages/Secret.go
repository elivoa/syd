package index

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service"
)

// _______________________________________________________________________________
//  ROOT Page
//
type Secret struct {
	core.Page
	UserToken *model.UserToken

	Current   int `path-param:"1"` // pager: the current item. in pager.
	PageItems int `path-param:"2"` // pager: page size.
}

func (p *Secret) SetupRender() *exit.Exit {
	userToken, err := service.Auth.Auth(p.W, p.R, model.ROLE_LOGIN|model.ROLE_Secret)
	if err != nil {
		return exit.Error(err)
	}
	p.UserToken = userToken
	return nil
}
