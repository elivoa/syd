package service

import (
	"errors"
	"fmt"
)

func validate_orderby(value string, candidate ...string) error {
	if candidate != nil && len(candidate) > 0 {
		if !_validate_candidates(value, candidate...) {
			return errors.New(fmt.Sprintf("Invalid `OrderBy` value, candidates:%v", candidate))
		}
	}
	// TODO: common checkout.
	return nil
}

func validate_order(value string) error {
	if !_validate_candidates(value, "asc", "desc") {
		return errors.New(fmt.Sprintf("Invalid `Order` value, candidates:`asc`,`desc`."))
	}
	return nil
}

func _validate_candidates(value string, candidate ...string) bool {
	for _, c := range candidate {
		if c == value {
			return true
		}
	}
	return false
}
