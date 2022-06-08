package mock

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"math/rand"
)

func ContactInfo(count int) []entity.ContactInfo {
	c := make([]entity.ContactInfo, count)
	for i := range c {
		c[i] = entity.ContactInfo{
			Name:        []string{"Jadon Haines", "Cadi Butt", "Evie-Mae Lin", "Nia Parra"}[rand.Intn(3)],
			PhoneNumber: fmt.Sprintf("09%v%v%v%v%v%v%v%v", rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9)),
			Address:     []string{"10", "20", "30", "99"}[rand.Intn(3)],
			District:    "คลองมหานาค",
			State:       "ป้อมปราบศัตรูพ่าย",
			Province:    "กรุงเทพมหานคร",
			Postcode:    "10100",
		}
	}
	return c
}
