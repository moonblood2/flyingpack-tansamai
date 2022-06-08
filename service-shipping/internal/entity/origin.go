package entity

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Origin struct {
	Id int64 `json:"id,omitempty"`
	ContactInfo
	CreateDelete
}

func (o Origin) Validate() error {
	return validation.ValidateStruct(&o,
		//Id
		validation.Field(&o.Id, validation.When(o.Id != 0, validation.Required)),
		//ContactInfo, let's ozzo call Validate().
		validation.Field(&o.ContactInfo),
	)
}
