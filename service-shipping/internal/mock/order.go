package mock

import (
	"fmt"
	"github.com/jna-distribution/service-shipping/internal/entity"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

func ProviderCode(count int) []entity.ProviderCode {
	rand.Seed(time.Now().Unix())
	c := make([]entity.ProviderCode, count)
	for i := range c {
		c[i] = []entity.ProviderCode{entity.ProviderCodeShippop, entity.ProviderCodeFlash}[rand.Intn(1)]
	}
	return c
}

func PaymentMethod(count int) []entity.OrderPaymentMethod {
	p := make([]entity.OrderPaymentMethod, count)
	for i := range p {
		p[i] = []entity.OrderPaymentMethod{entity.OrderPaymentMethod1, entity.OrderPaymentMethod2, entity.OrderPaymentMethod3, entity.OrderPaymentMethod4}[rand.Intn(3)]
	}
	return p
}

func OrderProduct(count int, userId string) []entity.OrderProduct {
	rand.Seed(time.Now().Unix())
	p := make([]entity.OrderProduct, count)
	for i := 0; i < count; i++ {
		p[i] = entity.OrderProduct{
			UserId:    userId,                  //55 - 60
			SenderId:  int64(1 + rand.Intn(3)), //1 - 3
			ProductId: int32(1 + rand.Intn(1)), //1 - 2
			Quantity:  int32(1 + rand.Intn(9)), //1 - 10
			PaymentMethod: []entity.OrderPaymentMethod{
				entity.OrderPaymentMethod1,
				entity.OrderPaymentMethod2,
				entity.OrderPaymentMethod3,
				entity.OrderPaymentMethod4,
			}[rand.Intn(3)],
		}
	}
	return p
}

func OrderParcel(count int, userId string) []entity.OrderParcel {
	rand.Seed(time.Now().Unix())
	o := make([]entity.OrderParcel, count)
	for i := range o {
		o[i].OrderParcelId = uuid.NewV4().String()
		o[i].UserId = userId
		o[i].SenderId = 2
		o[i].OriginId = 1
		o[i].DestinationId = 1
		o[i].ProviderCode = []entity.ProviderCode{entity.ProviderCodeShippop, entity.ProviderCodeFlash}[rand.Intn(1)]
		o[i].Price = float64(25 + rand.Intn(200-25))
		o[i].PaymentMethod = []entity.OrderPaymentMethod{entity.OrderPaymentMethod1, entity.OrderPaymentMethod2, entity.OrderPaymentMethod3, entity.OrderPaymentMethod4}[rand.Intn(3)]
		o[i].Weight = float32(100 + rand.Intn(5000-100))
		o[i].Width = float32(1 + rand.Intn(280-1))
		o[i].Length = float32(1 + rand.Intn(280-1))
		o[i].Height = float32(1 + rand.Intn(280-1))
	}
	return o
}

func OrderParcelShippop(count int, orderParcelId []string) []entity.OrderParcelShippop {
	rand.Seed(time.Now().Unix())
	o := make([]entity.OrderParcelShippop, count)
	for i := range o {
		o[i].OrderParcelId = orderParcelId[i]
		o[i].PurchaseId = int64(10000 + rand.Intn(999999-10000))
		o[i].Status = ""
		o[i].CourierCode = "FLE"
		o[i].CourierTrackingCode = ""
		o[i].TrackingCode = ""
		o[i].CODAmount = float64(100 + rand.Intn(2000-100))
		o[i].OrderParcelShippopFlash = entity.OrderParcelShippopFlash{
			ShippopFlashSortingCode: entity.ShippopFlashSortingCode{
				SortCode:        "16B-16265-01",
				DstCode:         "SSN_SP-สามเสนนอก",
				SortingLineCode: "C13",
			},
		}
	}
	return o
}

func OrderParcelFlash(count int, orderParcelId []string) []entity.OrderParcelFlash {
	o := make([]entity.OrderParcelFlash, count)
	for i := range o {
		b := make([]byte, 10)
		pno, _ := rand.Read(b)

		o[i].OrderParcelId = orderParcelId[i]
		o[i].Pno = fmt.Sprintf("%X", pno)
		o[i].State = int16(1 + rand.Intn(6))
		o[i].StateText = "-"
		o[i].CODAmount = float64(100 + rand.Intn(2000-100))
	}
	return o
}
