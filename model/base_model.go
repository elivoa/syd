package model

import (
	"strings"
)

type JsonResponse struct {
	Data interface{} `json:"data,omitempty"`

	// messages.
	Total   int `json:"total"`
	Items   int `json:"items"`
	Current int `json:"page"`

	Error string `json:"error,omitempty"` // real identification, tracking number
}

func NewJsonResponse(data interface{}) *JsonResponse {
	return &JsonResponse{Data: data}
}

func NewJsonErrorResponse(err string) *JsonResponse {
	return &JsonResponse{Error: err}
}

//
//
//
type Params map[string]interface{}

func (p Params) IsTrue(key string) bool {
	if value, ok := p[key]; ok {
		switch value.(type) {
		case bool:
			return value.(bool)
		case string:
			return strings.ToLower(value.(string)) == "true"
		case int:
			return value.(int) > 0
		}
	}
	return false
}

func (p Params) Int(key string) int {
	if value, ok := p[key]; ok {
		return value.(int)
	}
	return 0
}
