package valueobject

import (
	"errors"
)

var ErrEmptyProductName = errors.New("product name is empty")

type ProductName struct {
	value string
}

func NewProductName(value string) (ProductName, error) {
	if err := validateProductName(value); err != nil {
		return ProductName{}, err
	}
	return ProductName{value: value}, nil
}

func (p ProductName) Value() string {
	return p.value
}

func validateProductName(value string) error {
	if len(value) == 0 {
		return ErrEmptyProductName
	}
	return nil
}
