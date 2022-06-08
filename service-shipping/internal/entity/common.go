package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type ContactInfo struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	District    string `json:"district"` //District is sub unit of State.
	State       string `json:"state"`    //State is sub unit of Province.
	Province    string `json:"province"`
	Postcode    string `json:"postcode"`
}

type CreateDelete struct {
	CreatedAt time.Time `json:"created_at"`
	DeletedAt int64     `json:"-"`
}

func (c ContactInfo) Validate() error {
	return validation.ValidateStruct(&c,
		//Name, min 6, max 64.
		validation.Field(&c.Name, validation.Required, validation.Length(6, 256)),
		//PhoneNumber, min 8, max 10.
		validation.Field(&c.PhoneNumber, validation.Required, validation.Length(8, 10)),
		//Address, min 1, max 256.
		validation.Field(&c.Address, validation.Required, validation.Length(1, 256)),
		//District, min 1, max 128.
		validation.Field(&c.District, validation.Required, validation.Length(1, 128)),
		//State, min 1, max 128.
		validation.Field(&c.State, validation.Required, validation.Length(1, 128)),
		//Province, min 1, max 128.
		validation.Field(&c.Province, validation.Required, validation.Length(1, 128)),
		//Postcode, min 8, max 8.
		validation.Field(&c.Postcode, validation.Required, validation.Length(5, 5)),
	)
}
