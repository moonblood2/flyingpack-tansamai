package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	SenderType int16

	Sender struct {
		Id               int64      `json:"id"`
		SenderType       SenderType `json:"sender_type"`
		NationalIdNumber string     `json:"national_id_number"`
		TaxIdNumber      string     `json:"tax_id_number"`
		PassportNumber   string     `json:"passport_number"`
		BirthDate        string     `json:"birth_date"`
		ContactInfo
		CreateDelete
	}
)

const (
	SenderTypeNaturalPerson  SenderType = 1
	SenderTypeJuristicPerson SenderType = 2
)

func (s Sender) Validate() error {
	return validation.ValidateStruct(&s,
		//Id,
		validation.Field(&s.Id, validation.When(s.Id != 0, validation.Required)),
		//SenderType.
		validation.Field(&s.SenderType, validation.Required, validation.In(SenderTypeNaturalPerson, SenderTypeJuristicPerson)),
		//NationalIdNumber, require when PassportNumber is blank.
		validation.Field(&s.NationalIdNumber, validation.When(s.PassportNumber == "", validation.Required, validation.Length(13, 13))),
		//TaxIdNumber,
		validation.Field(&s.TaxIdNumber, validation.When(s.TaxIdNumber != "", validation.Required)),
		//PassportNumber, require when NationalIdNumber is blank.
		validation.Field(&s.PassportNumber, validation.When(s.NationalIdNumber == "", validation.Required)),
		//Date use RFC3339 that alias of ISO8610 for Date type in postgres. Postgres Date type use yyyy-mm-dd format.
		validation.Field(&s.BirthDate, validation.Required, validation.Date("2006-01-02")),
		//ContactInfo, call Validate() of ContactInfo.
		validation.Field(&s.ContactInfo),
	)
}
