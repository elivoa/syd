package sale

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service"
	"github.com/elivoa/syd/tools"
)

type ProductApi struct {
	core.Page
	UserToken *model.UserToken

	// input
	Id        int    `query:"id"`
	Tab       string `query:"tab"`
	Current   int    `query:"page"`  // pager: the current item. in pager.
	PageItems int    `query:"items"` // pager: page size.
}

func (p *ProductApi) Activate() {
	// allow cross domain access.
	tools.HttpAllowCrossDomainAccess(p.W)
	// TODO Other light-weight validation.
}

func (p *ProductApi) SetupRender() *exit.Exit {
	return nil
}

// Required params: tab, current, pageItems
func (p *ProductApi) Onlist() *exit.Exit {
	// p.fixPagerParameters()

	// process Tab.
	products, total, page, items, err := service.Product.GetProducts(model.Params{
		"first-letter": p.Tab,
		"current":      p.Current, // current page
		"page_items":   p.PageItems,
		"return_total": true,
	}, service.WITH_PRODUCT_DETAIL|service.WITH_PRODUCT_INVENTORY)
	if err != nil {
		// TODO return error
		panic(err)
	}

	// filter out values.
	filteroutProductListJson(products)

	resp := model.NewJsonResponse(products)
	resp.Total = total
	resp.Current = page
	resp.Items = items
	return exit.MarshalJson(resp)
}

// Required params: id
func (p *ProductApi) Ongset() *exit.Exit {
	// p.fixPagerParameters()

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

func (p *ProductApi) fixPagerParameters() {
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE // TODO default pager number. Config this.
	}
	if p.Current <= 0 {
		p.Current = 0
	}
}

// remove unnecessary values to generate json.
func filteroutProductListJson(products []*model.Product) {
	if products != nil && len(products) > 0 {
		//		empty_order_detail := []*model.OrderDetail{}
		for _, m := range products {
			if nil != m {
				m.Brand = ""
				m.Supplier = 0
				m.FactoryPrice = 0
				m.ShelfNo = ""
				m.Capital = ""
				m.CreateTime = nil
				m.UpdateTime = nil
			}
		}
	}

}
