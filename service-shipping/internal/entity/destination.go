package entity

import "github.com/go-ozzo/ozzo-validation/v4"

type Destination struct {
	Id int64 `json:"id,omitempty"`
	ContactInfo
	CreateDelete
}

func (d Destination) Validate() error {
	return validation.ValidateStruct(&d,
		//Id
		validation.Field(&d.Id, validation.When(d.Id != 0, validation.Required)),
		//ContactInfo, let's ozzo call Validate().
		validation.Field(&d.ContactInfo),
	)
}
