package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	ParcelShape struct {
		Weight float32 `json:"weight"`
		Width  float32 `json:"width"`
		Length float32 `json:"length"`
		Height float32 `json:"height"`
	}
	//Parcel use for save to database, MakeOrder request and call provider CreateOrder.
	Parcel struct {
		Key          int
		Id           uint32       `json:"id,omitempty"`
		Price        float64      `json:"price,omitempty"`         //Price in ThaiBath(฿) unit, use float64 because have Sa-Tang sub unit of Bath, and it's precise more than float32.
		ProviderCode ProviderCode `json:"provider_code,omitempty"` //ProviderCode is code for each provider.
		CourierCode  string       `json:"courier_code,omitempty"`  //CourierCode is sub code of ProviderCode, it use in Shippop.
		EnableCOD    bool         `json:"enable_cod,omitempty"`    //EnableCOD flag that tell this parcel will use COD service.
		CODAmount    float64      `json:"cod_amount,omitempty"`    //CODAmount use when EnableCOD is true.
		Origin       Origin       `json:"origin"`
		Destination  Destination  `json:"destination"`
		ParcelShape  `json:"parcel_shape"`
	}
)

func (p Parcel) Validate() error {
	return validation.ValidateStruct(&p,
		//Id
		validation.Field(&p.Id, validation.When(p.Id != 0, validation.Required)),
		//Price
		validation.Field(&p.Price, validation.When(p.Price != 0, validation.Required)),
		//ProviderCode, require.
		validation.Field(&p.ProviderCode, validation.Required, validation.In(ProviderCodeShippop, ProviderCodeFlash)),
		//CourierCode, validate when ProviderCode is Shippop.
		validation.Field(&p.CourierCode, validation.When(p.ProviderCode == ProviderCodeShippop, validation.Required)),
		//CODAmount, require when EnableCOD is true.
		validation.Field(&p.CODAmount, validation.When(p.EnableCOD == true, validation.Required)),
		//Origin, let's ozzo call Validate().
		validation.Field(&p.Origin, validation.Required),
		//Destination, let's ozzo call Validate().
		validation.Field(&p.Destination, validation.Required),
		//ParcelShape, let's ozzo call Validate().
		validation.Field(&p.ParcelShape, validation.Required),
	)
}

func (c ParcelShape) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Weight, validation.Required),
		validation.Field(&c.Width, validation.Required),
		validation.Field(&c.Length, validation.Required),
		validation.Field(&c.Height, validation.Required),
	)
}
