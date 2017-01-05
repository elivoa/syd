package service

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/db"
	"github.com/elivoa/syd/dal/persondao"
	"github.com/elivoa/syd/model"
)

type PersonService struct{}

func (s *PersonService) EntityManager() *db.Entity {
	return persondao.EntityManager()
}

// Get a list of persons.
//
// params: {
//   "type":          -- person type:'customer', 'factory'
//   "current":       -- 分页用，当前页
//   "page_items":    -- 分页用，每页条目数
//   "return_total":  -- 是否需要计算满足条件的所有数据数量。
//   "orderby":       -- order by "field"
//   "order":         -- asc/desc
// }
// withs: -- 数据集中还需要返回的数量.
// TODO make this fucntion powerful.
func (s *PersonService) GetPersons(params model.Params, withs Withs) (
	persons []*model.Person, total, page, items int, err error) {

	// 1. create queryparser.
	var parser = s.EntityManager().NewQueryParser()
	parser.Where()

	// 2. append where conditions.
	personType := params.String("type")
	if personType != "" {
		parser.And("type", personType)
	}

	// 3. get total if required.
	if params.IsTrue("return_total") {
		total, err = parser.Count()
		if err != nil {
			panic(err.Error())
		}
	}

	// 4. order or something else.
	orderby := params.String("orderby")
	if orderby != "" {
		if err = validate_orderby(orderby, "id", "createtime", "name"); err != nil {
			return
		}
		order := params.String("order")
		if err = validate_order(order); err != nil {
			return
		}
		parser.OrderBy(orderby, order)
	}

	// 5. add limit/pager info.
	page = params.Int("current")
	items = params.Int("page_items")
	if page <= 0 {
		page = 0
	}
	if items <= 0 {
		items = config.LIST_PRODUCT_SIZE_WITH_LETTER
	}
	parser.Limit(page*items, items) // pager

	// 6. fetch data.
	persons, err = s.List(parser, withs)
	return
}

func (s *PersonService) List(parser *db.QueryParser, withs Withs) ([]*model.Person, error) {
	if models, err := persondao.List(parser); err != nil {
		return nil, err
	} else {
		return models, nil
	}
}

func (s *PersonService) Get(field string, value interface{}) (*model.Person, error) {
	return persondao.Get(field, value)
}

func (s *PersonService) GetPersonById(id int) (*model.Person, error) {
	return s.Get(s.EntityManager().PK, id)
}

// return list of person
// func (s *PersonService) GetPersons(t person.Type) ([]*model.Person, error) {
// 	return persondao.ListAll(string(t))
// }

// --------------------------------------------------------------------------------
// The following is helper function to fill user to models.
func (s *PersonService) _batchFetchPerson(ids []int64) (map[int64]*model.Person, error) {
	return persondao.ListPersonByIdSet(ids...)
}

func (s *PersonService) BatchFetchPerson(ids ...int64) (map[int64]*model.Person, error) {
	return s._batchFetchPerson(ids)
}

func (s *PersonService) BatchFetchPersonByIdMap(idset map[int64]bool) (map[int64]*model.Person, error) {
	var idarray = []int64{}
	if idset != nil {
		for id, _ := range idset {
			idarray = append(idarray, id)
		}
	}
	return s._batchFetchPerson(idarray)
}
