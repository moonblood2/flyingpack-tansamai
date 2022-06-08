package mock

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"math/rand"
	"time"
)

func SenderType(count int) []entity.SenderType {
	rand.Seed(time.Now().Unix())
	t := make([]entity.SenderType, count)
	for i := range t {
		t[i] = []entity.SenderType{entity.SenderTypeNaturalPerson, entity.SenderTypeJuristicPerson}[rand.Intn(1)]
	}
	return t
}

func Sender(count int) []entity.Sender {
	s := make([]entity.Sender, count)
	for i := range s {
		s[i] = entity.Sender{
			SenderType:       SenderType(1)[0],
			NationalIdNumber: NationalId(),
			TaxIdNumber:      "",
			PassportNumber:   "",
			BirthDate:        fmt.Sprintf("19%v%v-%v%v-%v%v", Rand(7, 9), Rand(0, 9), Rand(0, 1), Rand(1, 9), Rand(0, 2), Rand(1, 9)),
			ContactInfo:      ContactInfo(1)[0],
		}
	}
	return s
}
