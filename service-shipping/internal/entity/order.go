package entity

type (
	OrderType          int16
	OrderPaymentMethod int16

	OrderParcel struct {
		OrderParcelId string `json:"order_parcel_id"`
		UserId        string `json:"user_id"`
		SenderId      int64  `json:"sender_id"`
		OriginId      int64  `json:"origin_id"`
		DestinationId int64  `json:"destination_id"`

		ProviderCode ProviderCode `json:"provider_code"`

		Price         float64            `json:"price"`
		PaymentMethod OrderPaymentMethod `json:"payment_method"`

		TrackingCode string  `json:"tracking_code"`
		Status       string  `json:"status"`
		CodAmount    float64 `json:"cod_amount"`

		ParcelShape

		CreateDelete
	}
	OrderProduct struct {
		Id        int64  `json:"id"`
		UserId    string `json:"user_id"`
		SenderId  int64  `json:"sender_id"`
		ProductId int32  `json:"product_id"`

		Quantity      int32              `json:"quantity"`
		PaymentMethod OrderPaymentMethod `json:"payment_method"`

		CreateDelete
	}
)

//order_parcel of each provider.
type (
	OrderParcelShippop struct {
		Id            int64  `json:"id"`
		OrderParcelId string `json:"order_parcel_id"`

		PurchaseId          int64   `json:"purchase_id"` // purchase_id from provider booking api.
		Status              string  `json:"status"`
		CourierCode         string  `json:"courier_code"`
		CourierTrackingCode string  `json:"courier_tracking_code"`
		TrackingCode        string  `json:"tracking_code"`
		CODAmount           float64 `json:"cod_amount"`

		OrderParcelShippopFlash

		CreateDelete
	}

	OrderParcelShippopFlash struct {
		Id                   int64 `json:"id"`
		OrderParcelShippopId int64 `json:"order_parcel_shippop_id"`
		ShippopFlashSortingCode
		CreateDelete
	}

	ShippopFlashSortingCode struct {
		SortCode        string `json:"sort_code"`
		DstCode         string `json:"dst_code"`
		SortingLineCode string `json:"sorting_line_code"`
	}

	OrderParcelFlash struct {
		Id            int64
		OrderParcelId string

		Pno       string
		State     int16
		StateText string
		CODAmount float64

		CreateDelete
	}
)

const (
	OrderPaymentMethod1 OrderPaymentMethod = 1
	OrderPaymentMethod2 OrderPaymentMethod = 2
	OrderPaymentMethod3 OrderPaymentMethod = 3
	OrderPaymentMethod4 OrderPaymentMethod = 4
)
