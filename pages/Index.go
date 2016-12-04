package index

import (
	"elivoa/carfilm/model"
	"elivoa/carfilm/service"
	"github.com/elivoa/got/core"
)

// _______________________________________________________________________________
//  ROOT Page
//
type Index struct {
	core.Page
	UserToken *model.UserToken

	Current   int `path-param:"1"` // pager: the current item. in pager.
	PageItems int `path-param:"2"` // pager: page size.
}

func (p *Index) SetupRender() {
	p.UserToken = service.User.RequireLogin(p.W, p.R)

}

func (p *Index) UrlTemplate() string {
	return "/{{Start}}/{{PageItems}}"
}
