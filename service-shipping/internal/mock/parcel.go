package mock

import (
	"github.com/jna-distribution/service-shipping/internal/entity"
	"math/rand"
	"time"
)

func ParcelShape(count int64) []entity.ParcelShape {
	rand.Seed(time.Now().Unix())
	s := make([]entity.ParcelShape, count)
	for i := range s {
		s[i].Weight = float32(1 + rand.Intn(50-1))
		s[i].Width = float32(1 + rand.Intn(50-1))
		s[i].Length = float32(1 + rand.Intn(50-1))
		s[i].Height = float32(1 + rand.Intn(50-1))
	}
	return s
}

func Parcel(count int) []entity.Parcel {
	parcels := make([]entity.Parcel, count)
	for i := range parcels {
		parcels[i] = entity.Parcel{
			ProviderCode: ProviderCode(1)[0], //Shippop
			CourierCode:  "FLE",
			EnableCOD:    false,
			Origin:       Origin(1)[0],
			Destination:  Destination(1)[0],
			ParcelShape:  ParcelShape(1)[0],
		}
	}
	return parcels
}
