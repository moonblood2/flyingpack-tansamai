package entity

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Product struct {
	Id     int32   `json:"id"`
	UserId string  `json:"user_id,omitempty"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"` //Price Thai bath unit(฿).
	CreateDelete
}

var (
	ProductNameRule = []validation.Rule{
		validation.Length(3, 64),
	}
	ProductPriceRule = []validation.Rule{
		validation.Min(1.0),
	}
)

func (p Product) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserId, validation.Required),
		validation.Field(&p.Name, ProductNameRule...),
		validation.Field(&p.Price, ProductPriceRule...),
	)
}
