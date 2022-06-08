package ordparcel

import "github.com/jna-distribution/service-shipping/internal/entity"

type FindRecord struct {
	OrderParcel             entity.OrderParcel              `json:"order_parcel"`
	Sender                  entity.Sender                   `json:"sender"`
	Origin                  entity.Origin                   `json:"origin"`
	Destination             entity.Destination              `json:"destination"`
	OrderParcelShippop      *entity.OrderParcelShippop      `json:"order_parcel_shippop,omitempty"`
	OrderParcelShippopFlash *entity.OrderParcelShippopFlash `json:"order_parcel_shippop_flash,omitempty"`
}

type Repository interface {
	Insert(orders []entity.OrderParcel) ([]entity.OrderParcel, error)
	UpdateOrder(order entity.OrderParcel) error
	Find(userId string, startDate, endDate string) ([]FindRecord, error)
	UpdateStatus(trackingCode, status string, providerCode entity.ProviderCode) error
	FindUser(trackingCode string, providerCode entity.ProviderCode) (string, error)
	FindIdByTrackingCode(trackingCode string, providerCode entity.ProviderCode) (string, error)
	FindOrderByIds(userId string, ids []string) (map[string]FindRecord, error)
}
