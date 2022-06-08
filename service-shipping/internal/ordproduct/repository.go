package ordproduct

import "github.com/jna-distribution/service-shipping/internal/entity"

type FindRecord struct {
	OrderProduct entity.OrderProduct `json:"order_product"`
	Sender       entity.Sender       `json:"sender"`
	Product      entity.Product      `json:"product"`
}

type Repository interface {
	Insert(products []entity.OrderProduct) ([]entity.OrderProduct, error)
	Find(userId string, startDate, endDate string) ([]FindRecord, error)
}
