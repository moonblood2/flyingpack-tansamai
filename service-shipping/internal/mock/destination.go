package mock

import (
	"github.com/jna-distribution/service-shipping/internal/entity"
	"math/rand"
	"time"
)

func Destination(count int) []entity.Destination {
	rand.Seed(time.Now().Unix())
	d := make([]entity.Destination, count)
	for i := 0; i < count; i++ {
		d[i] = entity.Destination{
			ContactInfo: ContactInfo(1)[0],
		}
	}
	return d
}
