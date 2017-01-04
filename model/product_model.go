package model

import (
	"fmt"
	"github.com/elivoa/syd/base/product"
	"strings"
	"time"
)

// TODO Design:model, how to split model data and fields.
type Product struct {
	Id           int            `json:"id,omitempty"`         // id
	Name         string         `json:"name,omitempty"`       // product name
	ProductId    string         `json:"product_id,omitempty"` // 传说中的货号
	Status       product.Status `json:"status,omitempty"`     //
	Brand        string         `json:"brand,omitempty"`
	Price        float64        `json:"price,omitempty"`
	Supplier     int            `json:"supplier,omitempty"`
	FactoryPrice float64        `json:"factory_price,omitempty"`
	Stock        int            `json:"stock,omitempty"`    // 库存量 || not used again?
	ShelfNo      string         `json:"shelf_no,omitempty"` // 货架号
	Capital      string         `json:"capital,omitempty"`  // captical letter to quick access.
	Note         string         `json:"note,omitempty"`
	CreateTime   *time.Time     `json:"create_time,omitempty"`
	UpdateTime   *time.Time     `json:"update_time,omitempty"`

	Pictures string `json:"-"` // picture keys splited by ';' filenamne can't contain ';'

	// additional information, not in persistence
	Colors     []string            `json:",omitempty"` // stores in product_properties table.
	Sizes      []string            `json:",omitempty"`
	Properties map[string][]string `json:",omitempty"` // other properties // TODO

	// stock information. format: map[color__size]nstock
	// special values in stock table
	//   stock = -1 means this pair of combination doesn't exist.
	//   stock = -2 means the pair is deleted.(may be price is available)
	Stocks Stocks `json:",omitempty"` // map[string]int
}

// TODO make a new structure of stocks;
type Stocks map[string]map[string]int

func NewStocks() Stocks {
	return map[string]map[string]int{}
}

// Create default empty Product
func NewProduct() *Product {
	now := time.Now()
	return &Product{
		Colors: []string{"", "", ""},
		Sizes:  []string{"", "", ""},
		// Stocks: map[string]int{},
		CreateTime: &now,
	}
}

func (s Stocks) Set(color, size string, stock int) {
	sizes, ok := s[color]
	if !ok {
		sizes = map[string]int{}
		s[color] = sizes
	}
	sizes[size] = stock
}

func (s Stocks) Loop(callback func(color, size string, stock int)) {
	for color, sizes := range s {
		if sizes != nil {
			for size, stock := range sizes {
				callback(color, size, stock)
			}
		}
	}
}

func (s Stocks) Total() int {
	total := 0
	for _, sizes := range s {
		if sizes != nil {
			for _, stock := range sizes {
				total += stock
			}
		}
	}
	return total
}

// Stock Item
type ProductStockItem struct {
	Color string
	Size  string
	Stock int
}

func (s ProductStockItem) Key() string {
	return fmt.Sprintf("%s__%s", s.Color, s.Size)
}

// func (p *Product) TotalStock() int {
// 	if nil != p.Stocks && len(p.Stocks) > 0 {
// 		var totalstock = 0
// 		for _, s := range p.Stocks {
// 			totalstock += s.Stock
// 		}
// 		return totalstock
// 	}
// 	return 0
// }

func (p *Product) ClearValues() {
	if p != nil {
		p.ClearColors()
		p.ClearSizes()
	}
}

func (p *Product) ClearColors() {
	newColors := []string{}
	if p.Colors != nil {
		hasValue := false
		for _, color := range p.Colors {
			color = strings.Trim(color, " ")
			if color != "" {
				hasValue = true
				newColors = append(newColors, color)
			}
		}
		if !hasValue {
			newColors = append(newColors, "默认颜色")
		}
	}
	p.Colors = newColors
}

func (p *Product) ClearSizes() {
	newSizes := []string{}
	if p.Sizes != nil {
		hasValue := false
		for _, size := range p.Sizes {
			size = strings.Trim(size, " ")
			if size != "" {
				hasValue = true
				newSizes = append(newSizes, size)
			}
		}
		if !hasValue {
			newSizes = append(newSizes, "均码")
		}
	}
	p.Sizes = newSizes
}

func (p *Product) PictureKeys() []string {
	res := strings.Split(p.Pictures, ";")
	pks := []string{}
	for _, pk := range res {
		if strings.TrimSpace(pk) != "" {
			pks = append(pks, strings.TrimSpace(pk))
		}
	}
	return pks
}
