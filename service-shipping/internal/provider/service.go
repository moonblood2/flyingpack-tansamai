package provider

import (
	"errors"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"time"
)

type (
	//resultShippop
	resultShippop struct {
		PurchaseID          int64                          `json:"purchase_id"` // purchase_id from provider booking api.
		CourierCode         string                         `json:"courier_code"`
		CourierTrackingCode string                         `json:"courier_tracking_code"`
		TrackingCode        string                         `json:"tracking_code"`
		CODAmount           float64                        `json:"cod_amount"`
		FlashSortingCode    entity.ShippopFlashSortingCode `json:"flash_sorting_code"`
	}
	//resultFlash
	resultFlash struct {
		Pno                string
		MchID              string
		OutTradeNo         string
		SortCode           string
		DstStoreName       string
		SortingLineCode    string
		EarlyFlightEnabled bool
		PackEnabled        bool
		CODAmount          float64
	}
	//BookingOrderResult
	BookingOrderResult struct {
		ParcelKey     int
		Status        bool                `json:"status"`
		Message       string              `json:"message"`
		TrackingCode  string              `json:"tracking_code"`
		ProviderCode  entity.ProviderCode `json:"provider_code"`
		ResultShippop resultShippop       //ResultShippop is response from booking order API of Shippop.
		ResultFlash   resultFlash         //ResultFlash is response from create order API of Flash Express.
	}
	//GetPriceOutput use with GetPrice method.
	GetPriceResult struct {
		Status  bool    `json:"status"`
		Message string  `json:"message,omitempty"`
		Price   float64 `json:"price"`
	}
	//InsertOrderInput
	InsertOrderInput struct {
		OrderParcelId string
		BookingOrderResult
	}

	Provider interface {
		//BookingOrder request to provider for booking order.
		BookingOrder(parcel entity.Parcel) (BookingOrderResult, error)
		//BookingOrders
		BookingOrders(parcels []entity.Parcel) ([]BookingOrderResult, error)
		//ConfirmOrder request to provider for confirm order.
		ConfirmOrder(result *BookingOrderResult) error
		//ConfirmOrders
		ConfirmOrders(results []BookingOrderResult) error
		//CancelOrder request to provider for cancel order.
		CancelOrder(result BookingOrderResult) error
		//GetPrice get service price from provider.
		GetPrice(parcel entity.Parcel) (GetPriceResult, error)
		//InsertOrder insert orders to database of each provider.
		InsertOrder(input []InsertOrderInput) error
		//UpdateOrderStatus
		UpdateOrderStatus(trackingCode, status, codStatus string) (orderParcelId string, statusCompletedDate, codTransferredDate *time.Time, err error)
	}

	List interface {
		Add(providerCode entity.ProviderCode, provider Provider)
		Get(providerCode entity.ProviderCode) (Provider, error)
		GroupingProvider(parcels []entity.Parcel) (map[entity.ProviderCode][]entity.Parcel, error)
	}

	providerList struct {
		providers map[entity.ProviderCode]Provider
		count     int
	}
)

func NewProviderList(providers map[entity.ProviderCode]Provider) *providerList {
	return &providerList{
		providers: providers,
		count:     len(providers),
	}
}

func (p *providerList) Add(providerCode entity.ProviderCode, provider Provider) {
	p.providers[providerCode] = provider
}

func (p *providerList) Get(providerCode entity.ProviderCode) (Provider, error) {
	//If provider_code is not exists in map.
	if _, ok := p.providers[providerCode]; !ok {
		return nil, internal.ErrInternal{InternalError: errors.New("ProviderList.Get, providerCode not found. ")}
	}
	return p.providers[providerCode], nil
}

//GroupingProvider group parcel for each shipping provider by provider code.
func (p *providerList) GroupingProvider(parcels []entity.Parcel) (map[entity.ProviderCode][]entity.Parcel, error) {
	parcelGroup := make(map[entity.ProviderCode][]entity.Parcel, p.count)
	//For each parcels.
	//TODO Validate provider_code for each parcels less than len(input.Parcels) more than or equal zero.
	for _, v := range parcels {
		parcelGroup[v.ProviderCode] = append(parcelGroup[v.ProviderCode], v)
	}
	return parcelGroup, nil
}
