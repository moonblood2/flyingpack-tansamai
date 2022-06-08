package mock

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"math/rand"
	"time"
)

func Origin(count int) []entity.Origin {
	rand.Seed(time.Now().Unix())
	o := make([]entity.Origin, count)
	for i := 0; i < count; i++ {
		o[i] = entity.Origin{
			ContactInfo: entity.ContactInfo{
				Name:        []string{"Jadon Haines", "Cadi Butt", "Evie-Mae Lin", "Nia Parra"}[rand.Intn(3)],
				PhoneNumber: fmt.Sprintf("09%v%v%v%v%v%v%v%v", rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9)),
				Address:     []string{"10", "20", "30", "99"}[rand.Intn(3)],
				District:    "คลองมหานาค",
				State:       "ป้อมปราบศัตรูพ่าย",
				Province:    "กรุงเทพมหานคร",
				Postcode:    "10100",
			},
		}
	}
	return o
}
