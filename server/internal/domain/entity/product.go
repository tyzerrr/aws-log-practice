package entity

import (
	"github.com/google/uuid"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/valueobject"
)

type Product struct {
	ID          uuid.UUID
	Name        valueobject.ProductName
	Description string
	PriceAmount valueobject.PriceAmount
	Currency    valueobject.Currency
	Archived    bool
}

func NewProduct(
	id uuid.UUID,
	name valueobject.ProductName,
	description string,
	price valueobject.PriceAmount,
	currency valueobject.Currency,
	archived bool,
) *Product {
	return &Product{
		ID:          id,
		Name:        name,
		Description: description,
		PriceAmount: price,
		Currency:    currency,
		Archived:    archived,
	}
}
