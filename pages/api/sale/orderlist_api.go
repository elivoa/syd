package sale

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service"
	"github.com/elivoa/syd/tools"
)

type OrderApi struct {
	core.Page
	UserToken *model.UserToken

	// input
	Tab       string `query:"tab"`
	Current   int    `query:"page"`  // pager: the current item. in pager.
	PageItems int    `query:"items"` // pager: page size.

	// output ?
	Total int // pager: total items available
}

func (p *OrderApi) Activate() {
	// allow cross domain access.
	tools.HttpAllowCrossDomainAccess(p.W)
	// TODO Other light-weight validation.
}
func (p *OrderApi) SetupRender() *exit.Exit {
	return nil
}

//
func (p *OrderApi) Onlist() *exit.Exit {
	p.fixPagerParameters()

	// process Tab.
	var status = p.Tab
	if status == "all" {
		status = ""
	}

	orders, total, err := service.Order.GetOrders(model.Params{
		"status":       status,
		"current":      p.Current, // crrent page
		"page_items":   p.PageItems,
		"return_total": true,
	}, service.WITH_PERSON)
	if err != nil {
		// TODO return error
		panic(err)
	}

	// filter out values.
	filteroutOrderListJson(orders)

	resp := model.NewJsonResponse(orders)
	resp.Total = total
	resp.Current = p.Current
	resp.Items = p.PageItems
	return exit.MarshalJson(resp)
}

func (p *OrderApi) fixPagerParameters() {
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE // TODO default pager number. Config this.
	}
	if p.Current <= 0 {
		p.Current = 0
	}
}

// remove unnecessary values to generate json.
func filteroutOrderListJson(orders []*model.Order) {
	if orders != nil && len(orders) > 0 {
		empty_order_detail := []*model.OrderDetail{}
		for _, o := range orders {
			if nil != o {
				o.DeliveryTrackingNumber = ""
				o.Accumulated = 0
				o.Details = empty_order_detail
				o.UpdateTime = nil
				o.CloseTime = nil
				if nil != o.Customer {
					o.Customer.Type = ""
					o.Customer.Phone = ""
					o.Customer.City = ""
					o.Customer.Address = ""
					o.Customer.PostalCode = 0
					o.Customer.QQ = 0
					o.Customer.Website = ""
					o.Customer.Note = ""
					o.Customer.CreateTime = nil
					o.Customer.UpdateTime = nil
				}
			}
		}
	}

}
