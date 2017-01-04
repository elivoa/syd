package model

import (
	"time"
)

//
// core person model
//
type Person struct {
	Id         int    `json:"id,omitempty"`   // id // TODO: change to int64
	Name       string `json:"name,omitempty"` // pesron name
	Type       string `json:"type,omitempty"` // `enum(客户Customer|厂家Factory)` // person type
	Phone      string `json:"phone,omitempty"`
	City       string `json:"city,omitempty"`
	Address    string `json:"address,omitempty"`
	PostalCode int    `json:"postal_code,omitempty"`
	QQ         int    `json:"qq,omitempty"`
	Website    string `json:"website,omitempty"`
	Note       string `json:"note,omitempty"`

	// Customer: 存储累计欠款; Factory: TODO
	AccountBallance float64 `json:"account_ballance,omitempty"`

	CreateTime *time.Time `json:"create_time,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`

	// TODO ++
	/* favorite delivery method: TakeAway|SFExpress|物流 */
	DeliveryMethod string `json:"delivery_method,omitempty"`

	// Fax        string
}

func NewPerson() *Person {
	return &Person{Name: "", Type: "customer", Note: ""}
}

func (p *Person) Accomulated() float64 {
	return -p.AccountBallance
}

func (p *Person) IsCustomer() bool {
	return p.Type == "customer"
}

func (p *Person) IsFactory() bool {
	return p.Type == "factory"
}

//
// Advanced Wrapper
//
type Customer struct {
	Person
	// advanced properties
	// Accumulated float64 // 累计欠款 // TODO replaced by AccountBallance
}

type Producer struct {
	Person
	// advanced properties
}

// TODO type is enum

//
// Customer Special Price
//

type CustomerPrice struct {
	Id           int
	PersonId     int
	ProductId    int
	Price        float64
	CreateTime   time.Time
	LastUsedTime time.Time
}
