package valueobject

import (
	"errors"
)

var ErrNegativeValue = errors.New("product price amount is negative value")

type PriceAmount struct {
	value int64
}

func NewProductAmount(value int64) (PriceAmount, error) {
	if err := validatePriceAmount(value); err != nil {
		return PriceAmount{}, err
	}
	return PriceAmount{value: value}, nil
}

func validatePriceAmount(value int64) error {
	if value < 0 {
		return ErrNegativeValue
	}
	return nil
}
