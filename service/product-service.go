package service

import (
	"fmt"
	"github.com/elivoa/got/db"
	"github.com/elivoa/syd/base/product"
	"path/filepath"
	"strings"
	// "github.com/elivoa/syd/dal/inventorydao"
	"github.com/elivoa/got/config"
	"github.com/elivoa/syd/dal/productdao"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service/suggest"
	"github.com/elivoa/syd/utils"
	"time"
)

type ProductService struct{}

func (s *ProductService) EntityManager() *db.Entity {
	return productdao.EntityManager()
}

// Get a list of orders.
//
// params: {
//   "first-letter":  -- order.status
//   "current":       -- 分页用，当前页
//   "page_items":    -- 分页用，每页条目数
//   "return_total": -- 是否需要计算满足条件的所有数据数量。
// }
// withs: -- 数据集中还需要返回的数量.
// TODO make this fucntion powerful.
func (s *ProductService) GetProducts(params model.Params, withs Withs) (
	products []*model.Product, total, page, items int, err error) {

	// fetch data
	var parser = s.EntityManager().NewQueryParser()
	parser.Where()

	firstLetter, ok := params["first-letter"]
	if ok {
		if firstLetter != "all" {
			parser.And("capital", firstLetter)
		}
	}

	// get total
	if params.IsTrue("return_total") {
		total, err = parser.Count()
		if err != nil {
			panic(err.Error())
		}
	}

	// 2. get order list.
	page = params.Int("current")
	items = params.Int("page_items")
	if page <= 0 {
		page = 0
	}
	if items <= 0 {
		if firstLetter == "all" || firstLetter == "" {
			items = config.LIST_PRODUCT_SIZE
		} else {
			items = config.LIST_PRODUCT_SIZE_WITH_LETTER
		}
	}
	parser.Limit(page*items, items) // pager
	products, err = s.List(parser, withs)
	return
}

//
// Get Product, with product's size and color properties.
// params: {
//   "id":            -- product id
// }
//
func (s *ProductService) GetProduct(id int /* params model.Params,*/, withs Withs) (*model.Product, error) {
	if product, err := productdao.Get(id); err != nil {
		return nil, err
	} else if nil != product {
		models := []*model.Product{product}
		if withs&WITH_PRODUCT_DETAIL > 0 {
			if err := productdao.FillProductPropertiesByIdSet(models); err != nil {
				return nil, err
			}
		}
		// if withs&WITH_PRODUCT_INVENTORY > 0 {
		// 	if err := inventorydao.FillProductStocksByIdSet(models); err != nil {
		// 		return nil, err
		// 	}
		// }
		return product, nil
	}
	return nil, nil
}

func (s *ProductService) GetFullProduct(id int) (*model.Product, error) {
	return s.GetProduct(id, WITH_PRODUCT_DETAIL|WITH_PRODUCT_INVENTORY)
}

//////////////////////////////////////////////////////////////////////////////////

//
// CreateProduct create a new Product into database, including it's properties.
//
func (s *ProductService) CreateProduct(product *model.Product) (*model.Product, error) {
	if product == nil {
		panic("Product can't be null!")
	}

	product.Capital = s.getCapital(product.Name)
	newProduct, err := productdao.Create(product)
	if err != nil {
		return nil, err
	}
	// newProduct := dal.CreateProduct(product)
	if product.Colors != nil {
		productdao.UpdateProductProperties(newProduct.Id, "color", product.Colors...)
	}
	if product.Sizes != nil {
		productdao.UpdateProductProperties(newProduct.Id, "size", product.Sizes...)
	}

	// update suggest
	suggest.Add(suggest.Product, newProduct.Name, newProduct.Id, product.ProductId)

	return newProduct, nil
}

// func (s *ProductService) UpdateProduct(product *model.Product) {
// 	if product == nil {
// 		return
// 	}
// 	// update product information
// 	product.Capital = s.getCapital(product.Name)
// 	if _, err := productdao.UpdateProduct(product); err != nil {
// 		panic(err.Error())
// 	}

// 	// update it's properties
// 	if product.Colors != nil {
// 		productdao.UpdateProductProperties(product.Id, "color", product.Colors...)
// 	}
// 	if product.Sizes != nil {
// 		productdao.UpdateProductProperties(product.Id, "size", product.Sizes...)
// 	}

// 	// update stock information
// 	if product.Stocks != nil {
// 		// TODO change to edit/create/delete;
// 		inventorydao.ClearProductStock(product.Id) // clear
// 		product.Stocks.Loop(func(color, size string, stock int) {
// 			inventorydao.SetProductStock(product.Id, color, size, stock)
// 		})
// 	}

// 	// update suggest
// 	fmt.Println("^^^^^^ update product, ", product)
// 	fmt.Println("Productid is: ", product.ProductId)
// 	suggest.Update(suggest.Product, product.Name, product.Id, product.ProductId)

// }

func (s *ProductService) DeleteProduct(id int) (affacted int64, err error) {
	if affacted, err = productdao.Delete(id); err != nil {
		return -1, err
	} else {
		suggest.Delete(suggest.Product, id)
		return
	}
}

func (s *ProductService) ChangeStatus(id int, status product.Status) (affacted int64, err error) {
	if affacted, err = productdao.ChangeStatus(id, status); err != nil {
		return -1, err
	} else {
		// TODO affact suggest ???? should change status's status.
		// suggest.Delete(suggest.Product, id)
		return
	}
}

// --------------------------------------------------------------------------------
// The following is helper function to fill user to models.
func (s *ProductService) _batchFetchProduct(ids []int64) (map[int64]*model.Product, error) {
	return productdao.ListProductsByIdSet(ids...)
}

func (s *ProductService) BatchFetchProduct(ids ...int64) (map[int64]*model.Product, error) {
	return s._batchFetchProduct(ids)
}

func (s *ProductService) BatchFetchProductByIdMap(idset map[int64]bool) (map[int64]*model.Product, error) {
	var idarray = []int64{}
	if idset != nil {
		for id, _ := range idset {
			idarray = append(idarray, id)
		}
	}
	return s._batchFetchProduct(idarray)
}

func (s *ProductService) getCapital(text string) string {
	str := utils.ParsePinyin(text)
	if len(str) > 0 {
		return strings.ToLower(str[0:1])
	}
	return "-"
}

func (s *ProductService) List(parser *db.QueryParser, withs Withs) ([]*model.Product, error) {
	if models, err := productdao.List(parser); err != nil {
		return nil, err
	} else {
		// TODO: Print warrning information when has unused withs.
		// fmt.Println("--------------------------------------------------------------------", withs)
		if withs&WITH_PRODUCT_DETAIL > 0 {
			if err := productdao.FillProductPropertiesByIdSet(models); err != nil {
				return nil, err
			}
		}
		// if withs&WITH_PRODUCT_INVENTORY > 0 {
		// 	if err := inventorydao.FillProductStocksByIdSet(models); err != nil {
		// 		return nil, err
		// 	}
		// }
		return models, nil
	}
}

func (s *ProductService) ListStocks(parser *db.QueryParser) ([]*model.Product, error) {

	// func FillProductStocksByIdSet(models []*model.Product) error {
	// 	if nil == models || len(models) == 0 {
	// 		return nil
	// 	}

	// 	var idset = map[int64]bool{}
	// 	for _, m := range models {
	// 		idset[int64(m.Id)] = true
	// 	}
	// 	if allstocks, err := GetAllStocksByIdSet(idset); err != nil {
	// 		return err
	// 	} else {
	// 		if nil != allstocks {
	// 			for _, m := range models {
	// 				if stock, ok := allstocks[int64(m.Id)]; ok {
	// 					m.Stocks = stock
	// 					m.Stock = stock.Total()
	// 				}
	// 			}
	// 		}
	// 	}
	// 	return nil
	// }

	if models, err := productdao.List(parser); err != nil {
		return nil, err
	} else {
		// if err := inventorydao.FillProductStocksByIdSet(models); err != nil {
		// 	return nil, err
		// }
		return models, nil
	}
}

// No use // TODO: delete this;
// func (s *ProductService) ListProductsByCapital(capital string) ([]*model.Product, error) {
// 	return productdao.ListByCapital(capital)
// }

// Non-standard fill.
// func (s *ProductService) FillProductsWithDetails(models []*model.Product) error {
// 	var idset = map[int64]bool{}
// 	for _, model := range models {
// 		idset[int64(model.Id)] = true
// 	}
// 	productdao.FillProductPropertiesByIdSet(models)

// 	personmap, err := Product.BatchFetchProductDetailsByIdMap(idset)
// 	if err != nil {
// 		return err
// 	}
// 	if nil != personmap && len(personmap) > 0 {
// 		for _, order := range orders {
// 			if person, ok := personmap[int64(order.CustomerId)]; ok {
// 				order.Customer = person
// 			}
// 		}
// 	}
// 	return nil
// }

// --------------------------------------------------------------------------------
func (s *ProductService) RebuildProductPinyinCapital() {
	fmt.Println("________________________________________________________________________________")
	fmt.Println("Rebuild Product Capital")

	qp := db.NewQueryParser().Limit(10000).Where()
	products, err := s.List(qp, 0)
	if err != nil {
		panic(err.Error())
	}
	for _, product := range products {
		product.Capital = s.getCapital(product.Name)
		if _, err := productdao.UpdateProduct(product); err != nil {
			panic(err.Error())
		}
		fmt.Printf("> processing %v capital is: %v\n", product.Name, product.Capital)
	}
	fmt.Println("all done")
}

func (s *ProductService) ProductPictrues(product *model.Product) []string {
	if nil == product {
		return []string{}
	}
	pkeys := product.PictureKeys()
	for i := 0; i < len(pkeys); i++ {
		pkeys[i] = filepath.Join("/pictures", pkeys[i])
	}
	return pkeys
}

func (p *ProductService) ProductPicture(product *model.Product, index int) string {
	pictures := p.ProductPictrues(product)
	if nil != pictures && len(pictures) > index {
		return pictures[index]
	}
	return ""
}

/*
 * Stat: StatDailySalesData - 统计产品每日销售数量
 * productId == 0 -- Read all day sales.
 */
// func (s *ProductService) StatDailySalesData(productId int, period, combine_days int, endday time.Time) (
// 	model.ProductSales, error) {

// 	dprint := false
// 	remove_year := true
// 	default_period := 30
// 	if combine_days == 0 { // 默认点合并策略
// 		switch period {
// 		case 7:
// 			combine_days = 1
// 		case 30:
// 			combine_days = 5
// 		case 90:
// 			combine_days = 7
// 		case 365:
// 			combine_days = 7
// 		case 1095:
// 			combine_days = 30
// 		default:
// 			combine_days = 1
// 		}
// 	}

// 	showdays := period
// 	if showdays <= 0 {
// 		showdays = default_period
// 		combine_days = 1
// 	}
// 	keys := datekeys(showdays, endday)
// 	if keys == nil || len(keys) <= 0 {
// 		return nil, nil
// 	}

// 	if salesdata, err := productdao.DailySalesData(productId, keys[0], true, endday); err != nil {
// 		panic(err)
// 	} else {
// 		newps := model.ProductSales{}
// 		for _, key := range keys {
// 			// performance issue?
// 			found := false
// 			for _, v := range salesdata {
// 				if v.Key == key {
// 					found = true
// 					newps = append(newps, v)
// 					break
// 				}
// 			}
// 			if !found {
// 				newps = append(newps, &model.SalesNode{Key: key, Value: 0})
// 			}
// 		}

// 		if remove_year { // remove year-
// 			for _, node := range newps {
// 				if node != nil && len(node.Key) > 5 {
// 					node.Key = node.Key[5:]
// 				}
// 			}
// 		}

// 		if dprint {
// 			fmt.Println("\nDEVELOP .................................................")
// 			for _, node := range newps {
// 				fmt.Println("\t", node.Key, " is ", node.Value)
// 			}
// 		}

// 		if combine_days > 1 {
// 			var (
// 				idx          int = 0
// 				first_key    string
// 				last_key     string
// 				current      *model.SalesNode
// 				combinedNode *model.SalesNode
// 				combinedps   = model.ProductSales{}
// 				start        = false
// 			)
// 			for i := len(newps) - 1; i >= 0; i-- {
// 				current = newps[i]

// 				switch idx % combine_days {
// 				case 0: // every first one
// 					// fmt.Println(" - start", i)
// 					start = true
// 					last_key = current.Key
// 					combinedNode = &model.SalesNode{}

// 				case combine_days - 1: // end
// 					first_key = current.Key
// 					combinedNode.Key = fmt.Sprintf("%s,%s", first_key, last_key)
// 					// combinedNode.Value = combinedNode.Value / combine_days
// 					combinedps = append(combinedps, combinedNode)
// 					// fmt.Println(" - end", i, ": combined is : ", combinedNode.Value)
// 					start = false

// 				default: // in middle
// 				}
// 				combinedNode.Value += current.Value // combine values.
// 				// fmt.Printf("idx: %d mod %d = %d\n", idx, combine_days, (idx % combine_days))
// 				idx += 1
// 			}

// 			if start { // the last one
// 				first_key = current.Key
// 				combinedNode.Key = fmt.Sprintf("%s - %s", first_key, last_key)
// 				combinedps = append(combinedps, combinedNode)
// 			}

// 			if dprint {
// 				fmt.Println("\nDEVELOP:: Combined ProductSales ....................................")
// 				for _, node := range combinedps {
// 					fmt.Println("\t", node.Key, " is ", node.Value)
// 				}
// 			}

// 			// 需要翻转数组
// 			ncps := model.ProductSales{}
// 			for i := len(combinedps) - 1; i >= 0; i-- {
// 				ncps = append(ncps, combinedps[i])
// 			}
// 			return ncps, nil
// 		}

// 		return newps, nil
// 	}
// }

// endday - 最后日期
// lastNDays - 从endday往前数N天。
func datekeys(lastNDays int, endday time.Time) []string {
	t := endday.AddDate(0, 0, -lastNDays+1)
	// t := time.Now().AddDate(0, 0, -lastNDays+1)
	result := []string{}
	for i := 0; i < lastNDays; i++ {
		result = append(result, t.AddDate(0, 0, i).Format("2006-01-02"))
	}

	return result
}
