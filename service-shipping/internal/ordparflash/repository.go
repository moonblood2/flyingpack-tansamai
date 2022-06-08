package ordparflash

import "github.com/jna-distribution/service-shipping/internal/entity"

type Repository interface {
	Insert(orders []entity.OrderParcelFlash) ([]entity.OrderParcelFlash, error)
}
