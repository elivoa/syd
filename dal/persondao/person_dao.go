/*
  Data Access Object for person module.

  Time-stamp: <[person_dao.go] Elivoa @ Thursday, 2017-01-05 13:19:55>

  Note: This is the latest Template for dao functions.

*/
package persondao

import (
	"database/sql"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/db"
	"github.com/elivoa/syd/model"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// debug option
var logdebug = true

var core_fields = []string{
	"name", "type", "phone", "city", "address", "postalcode", "qq", "website", "note",
	"AccountBallance", "createtime",
}

var em = &db.Entity{
	Table:        "person",
	PK:           "id",
	Fields:       append(append([]string{"id"}, core_fields...), "updatetime"),
	CreateFields: core_fields,
	UpdateFields: core_fields,
}

func EntityManager() *db.Entity {
	return em
}

func init() {
	db.RegisterEntity("person", em)
}

// ________________________________________________________________________________
// Get person by person id
//

// // new version
// func GetPersonById(id int) (*model.Person, error) {
// 	return Get(em.PK, id)
// }

// new version
func Get(field string, value interface{}) (*model.Person, error) {
	var query = em.Select().Where(field, value)
	return _one(query)
}

// the last part, read the list from rows
func _list(query *db.QueryParser) ([]*model.Person, error) {
	models := make([]*model.Person, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.Person{}
			err := rows.Scan(
				&m.Id, &m.Name, &m.Type, &m.Phone, &m.City, &m.Address, &m.PostalCode, &m.QQ,
				&m.Website, &m.Note, &m.AccountBallance, &m.CreateTime, &m.UpdateTime,
			)
			models = append(models, m)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return models, nil
}

// only return the first result;
func _one(query *db.QueryParser) (*model.Person, error) {
	m := &model.Person{}
	if err := query.Query( // TODO: change to QueryOne
		func(rows *sql.Rows) (bool, error) {
			err := rows.Scan(
				&m.Id, &m.Name, &m.Type, &m.Phone, &m.City, &m.Address, &m.PostalCode, &m.QQ,
				&m.Website, &m.Note, &m.AccountBallance, &m.CreateTime, &m.UpdateTime,
			)
			return false, err // don't fetch the second line. first is enough;
		},
	); err != nil {
		return nil, err
	}
	return m, nil
}

// personType: customer, factory
// The very old method. Change to new version use person.Type
func ListAll(personType string) ([]*model.Person, error) {
	persons := make([]*model.Person, 0)
	err := em.Select().Where("type", personType).Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.Person)
			err := rows.Scan(
				&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ,
				&p.Website, &p.Note, &p.AccountBallance, &p.CreateTime, &p.UpdateTime,
			)
			persons = append(persons, p)
			return true, err
		},
	)
	if err != nil {
		return nil, err
	}
	return persons, nil
}

// ________________________________________________________________________________
// Create person
//
func Create(person *model.Person) error {
	res, err := em.Insert().Exec(
		person.Name, person.Type, person.Phone, person.City, person.Address, person.PostalCode,
		person.QQ, person.Website, person.Note, person.AccountBallance, time.Now(),
		// TODO: later change to create time outside.
	)
	if err != nil {
		return err
	}
	liid, err := res.LastInsertId()
	person.Id = int(liid)
	return nil
}

func List(parser *db.QueryParser) ([]*model.Person, error) {
	parser.SetEntity(em) // set entity manager into query parser.
	parser.Reset()       // to prevent if parser is used before. TODO:Is this necessary?
	// append default behavore.
	parser.DefaultOrderBy("createtime", db.DESC)
	parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	parser.Select()
	return _list(parser)
}

// ________________________________________________________________________________
// Update returns RowsAffacted, error
//
func Update(person *model.Person) (int64, error) {
	res, err := em.Update().Exec(
		person.Name, person.Type, person.Phone, person.City, person.Address, person.PostalCode, person.QQ,
		person.Website, person.Note, person.AccountBallance, time.Now(),
		person.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

var EmptyPersonList = &[]model.Person{}

// ________________________________________________________________________________
// Delete a pesron
//
func Delete(id int) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func ListPersonByIdSet(ids ...int64) (map[int64]*model.Person, error) {
	if nil == ids || len(ids) == 0 {
		return nil, nil
	}
	var query *db.QueryParser
	parser := em.Select().Where()
	query = parser.InInt64(em.PK, ids...).OrderBy(em.PK, db.DESC)

	models := make([]*model.Person, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.Person)
			err := rows.Scan(
				&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ,
				&p.Website, &p.Note, &p.AccountBallance, &p.CreateTime, &p.UpdateTime,
			)
			models = append(models, p)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	// return the map
	var resultmap = map[int64]*model.Person{}
	for _, u := range models {
		resultmap[int64(u.Id)] = u
	}
	return resultmap, nil
}

// the last part, read the list from rows
// func _list(query *db.QueryParser) ([]*model.Person, error) {
// 	models := make([]*model.Person, 0)
// 	if err := query.Query(
// 		func(rows *sql.Rows) (bool, error) {
// 			p := &model.Person{}
// 			err := rows.Scan(
// 				&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ,
// 				&p.Website, &p.Note, &p.AccountBallance, &p.CreateTime, &p.UpdateTime)
// 			models = append(models, p)
// 			return true, err
// 		},
// 	); err != nil {
// 		return nil, err
// 	}
// 	return models, nil
// }
