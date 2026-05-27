package valueobject

import (
	"errors"
)

var ErrEmptyCurrency = errors.New("currency is empty")

type Currency struct {
	value string
}

func NewCurrency(value string) (Currency, error) {
	if err := validateCurrency(value); err != nil {
		return Currency{}, err
	}
	return Currency{value: value}, nil
}

func validateCurrency(value string) error {
	if len(value) == 0 {
		return ErrEmptyCurrency
	}
	return nil
}
