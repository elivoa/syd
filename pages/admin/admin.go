package admin

import (
	"github.com/elivoa/got/core"
	// "github.com/elivoa/syd/service"
)

//________________________________________________________________________________
//
//
type AdminIndex struct {
	core.Page
}

func (p *AdminIndex) OnRebuildProductPinyin() {
	// service.Product.RebuildProductPinyinCapital()
}
