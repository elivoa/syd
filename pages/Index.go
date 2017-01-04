package index

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service"
	"time"
)

type Index struct {
	core.Page
	UserToken *model.UserToken

	// Current   int `path-param:"1"` // pager: the current item. in pager.
	// PageItems int `path-param:"2"` // pager: page size.
}

func (p *Index) SetupRender() *exit.Exit {
	// p.UserToken = service.User.RequireLogin(p.W, p.R)
	var err error
	if p.UserToken, err = service.Auth.OAuthToken(p.W, p.R); err != nil {
		return exit.Error(err)
	}
	return nil
}

func (p *Index) Expired() bool {
	if p.UserToken != nil && p.UserToken.Token != nil {
		if time.Now().Before(p.UserToken.Token.Expiry) {
			return true
		}

	}
	return false
}
