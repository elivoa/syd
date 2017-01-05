package sale

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service"
	"github.com/elivoa/syd/tools"
)

type Customer struct {
	core.Page
	UserToken *model.UserToken

	// input parameters.
	Id        int    `query:"id"`
	Tab       string `query:"tab"`   // main category
	Current   int    `query:"page"`  // pager: the current item. in pager.
	PageItems int    `query:"items"` // pager: page size.

	OrderBy string `query:"orderby"`
	Order   string `query:"order"`
}

func (p *Customer) Activate() {
	// allow cross domain access.
	tools.HttpAllowCrossDomainAccess(p.W)
	// TODO Other light-weight validation.
}

func (p *Customer) SetupRender() *exit.Exit {
	return nil
}

// Required params: tab, current, pageItems
func (p *Customer) OnlistCustomer() *exit.Exit {
	return p.listPerson("customer")
}

func (p *Customer) OnlistFactory() *exit.Exit {
	return p.listPerson("factory")
}

func (p *Customer) listPerson(personType string) *exit.Exit {
	persons, total, page, items, err := service.Person.GetPersons(model.Params{
		// ignore tab, use as
		"type":         personType,
		"current":      p.Current, // current page
		"page_items":   p.PageItems,
		"return_total": true,
		"orderby":      p.OrderBy,
		"order":        p.Order,
	}, service.WITH_NONE)
	if err != nil {
		// TODO return error
		panic(err)
	}

	// filter out values.
	filterCustomerListJson(persons)

	resp := model.NewJsonResponse(persons)
	resp.Total = total
	resp.Current = page
	resp.Items = items
	return exit.MarshalJson(resp)
}

// Required params: id
func (p *Customer) Onget() *exit.Exit {
	// process Tab.
	product, err := service.Product.GetFullProduct(p.Id)
	// products, total, page, items, err := service.Product.GetProducts(model.Params{
	//	"first-letter": p.Tab,
	//	"current":      p.Current, // current page
	//	"page_items":   p.PageItems,
	//	"return_total": true,
	// }, service.WITH_PRODUCT_DETAIL|service.WITH_PRODUCT_INVENTORY)
	if err != nil {
		// TODO return error
		panic(err)
	}

	// filter out values.
	// filteroutProductListJson(products)

	// resp := model.NewJsonResponse(products)
	// resp.Total = total
	// resp.Current = page
	// resp.Items = items
	return exit.MarshalJson(product)
}

func (p *Customer) fixPagerParameters() {
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE // TODO default pager number. Config this.
	}
	if p.Current <= 0 {
		p.Current = 0
	}
}

// remove unnecessary values to generate json.
func filterCustomerListJson(persons []*model.Person) {
	if persons != nil && len(persons) > 0 {
		for _, m := range persons {
			if nil != m {
				m.UpdateTime = nil
			}
		}
	}

}
