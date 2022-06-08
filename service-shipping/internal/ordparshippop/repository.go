package ordparshippop

import (
	"github.com/jna-distribution/service-shipping/internal/entity"
	"time"
)

type Repository interface {
	Insert(orders []entity.OrderParcelShippop) ([]entity.OrderParcelShippop, error)
	UpdateOrderStatus(trackingCode, status, codStatus string) (orderParcelId string, statusCompletedDate, codTransferredDate *time.Time, err error)
}
