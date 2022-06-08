package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/destination"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/ordparcel"
	"github.com/jna-distribution/service-shipping/internal/ordproduct"
	"github.com/jna-distribution/service-shipping/internal/origin"
	"github.com/jna-distribution/service-shipping/internal/product"
	"github.com/jna-distribution/service-shipping/internal/provider"
	"github.com/jna-distribution/service-shipping/internal/sender"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

type Service interface {
	MakeOrder(input MakeOrderInput, userId string) (MakeOrderOutput, error)
	GetParcelPrice(input GetParcelPriceInput) (GetParcelPriceOutput, error)
	GetOrderParcel(input GetOrderParcelInput) (GetOrderParcelOutput, error)
	UpdateOrderParcel(input UpdateOrderParcelInput) error
	GetOrderProduct(input GetOrderProductInput) (GetOrderProductOutput, error)
	UpdateOrderParcelStatus(input StatusData, providerCode entity.ProviderCode) error
	GetOrderParcelByIds(input GetOrderParcelByIdsInput) (GetOrderParcelByIdsOutput, error)
}

type service struct {
	productRepository product.Repository

	senderRepository      sender.Repository
	originRepository      origin.Repository
	destinationRepository destination.Repository

	orderParcelRepository  ordparcel.Repository
	orderProductRepository ordproduct.Repository

	providers provider.List

	jnaAddress         entity.Origin
	jnaAddressAsOrigin bool

	anUrl string
}

func NewService(
	productRepository product.Repository,

	senderRepository sender.Repository,
	originRepository origin.Repository,
	destinationRepository destination.Repository,

	orderParcelRepository ordparcel.Repository,
	orderProductRepository ordproduct.Repository,

	providers provider.List,

	jnaAddress entity.Origin,
	jnaAddressAsOrigin bool,

	anUrl string,
) *service {
	return &service{
		productRepository: productRepository,

		senderRepository:      senderRepository,
		originRepository:      originRepository,
		destinationRepository: destinationRepository,

		orderParcelRepository:  orderParcelRepository,
		orderProductRepository: orderProductRepository,

		providers: providers,

		jnaAddress:         jnaAddress,
		jnaAddressAsOrigin: jnaAddressAsOrigin,

		anUrl: anUrl,
	}
}

type (
	//Make order
	Product struct {
		Id       int32 `json:"id"`
		Quantity int32 `json:"quantity"`
	}
	MakeOrderInput struct {
		Sender        entity.Sender             `json:"sender"`
		PaymentMethod entity.OrderPaymentMethod `json:"payment_method"`
		Parcels       []entity.Parcel           `json:"parcels"`
		AnParcels     []entity.AnParcel         `json:"an_parcels"`
		Products      []Product                 `json:"products"`
	}
	makeOrderOutputStatus struct {
		Status  bool   `json:"status"`
		Message string `json:"message,omitempty"`
	}
	makeOrderOutputParcel struct {
		makeOrderOutputStatus
		Price                   float64                         `json:"price,omitempty"`
		TrackingCode            string                          `json:"tracking_code"`
		ShippopFlashSortingCode *entity.ShippopFlashSortingCode `json:"shippop_flash_sorting_code,omitempty"`
	}
	makeOrderOutputAnParcel struct {
		makeOrderOutputStatus
		TrackingCode            string                          `json:"tracking_code"`
		OrderParcelId           string                          `json:"order_parcel_id"`
		ShippopFlashSortingCode *entity.ShippopFlashSortingCode `json:"shippop_flash_sorting_code,omitempty"`
	}
	MakeOrderOutput struct {
		Status    bool                      `json:"status"`
		Parcels   []makeOrderOutputParcel   `json:"parcels,omitempty"`
		AnParcels []makeOrderOutputAnParcel `json:"an_parcels,omitempty"`
		Products  []makeOrderOutputStatus   `json:"products,omitempty"`
	}

	//Get parcel price
	parcelPrice struct {
		Status  bool    `json:"status"`
		Message string  `json:"message,omitempty"`
		Price   float64 `json:"price"`
	}
	GetParcelPriceInput  []entity.Parcel
	GetParcelPriceOutput []parcelPrice

	//Get order parcel
	GetOrderParcelInput struct {
		UserId    string `json:"user_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	GetOrderParcelOutput []ordparcel.FindRecord

	//Get order parcel by Ids
	GetOrderParcelByIdsInput struct {
		Ids    []string `json:"ids"`
		UserId string   `json:"user_id"`
	}
	GetOrderParcelByIdsOutput map[string]ordparcel.FindRecord

	//Get order product
	GetOrderProductInput struct {
		UserId    string `json:"user_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	GetOrderProductOutput []ordproduct.FindRecord

	UpdateOrderParcelInput struct {
		ParcelType    uint8  `json:"parcel_type"`
		OrderParcelId string `json:"order_parcel_id"`
		TrackingCode  string `json:"tracking_code"`
	}

	//StatusData use at UpdateStatus
	StatusData struct {
		TrackingCode string
		Status       string
		CodStatus    string
		Weight       float64
		Price        float64
		Datetime     time.Time
		Width        float64
		Height       float64
		Length       float64
	}

	WeightHookPayload struct {
		SpOrderParcelId string  `json:"sp_order_parcel_id"`
		TrackingCode    string  `json:"tracking_code"`
		Weight          float64 `json:"weight"`
		Width           float64 `json:"width"`
		Length          float64 `json:"length"`
		Height          float64 `json:"height"`
		StatusHookPayload
	}
	StatusHookPayload struct {
		SpOrderParcelId     string     `json:"sp_order_parcel_id"`
		TrackingCode        string     `json:"tracking_code"`
		Status              string     `json:"status"`
		CodStatus           string     `json:"cod_status"`
		StatusCompletedDate *time.Time `json:"status_completed_date"`
		CodTransferredDate  *time.Time `json:"cod_transferred_date"`
	}
)

//Validate validate product.
func (p Product) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Id, validation.Required, validation.Min(1)),
		validation.Field(&p.Quantity, validation.Required, validation.Min(1)),
	)
}

//Validate validate input from make order API.
func (i MakeOrderInput) Validate() error {
	var rules []*validation.FieldRules
	if len(i.Parcels) > 0 || len(i.Products) > 0 {
		rules = append(rules,
			validation.Field(&i.Sender),
			validation.Field(&i.PaymentMethod, validation.Required, validation.In(
				entity.OrderPaymentMethod1,
				entity.OrderPaymentMethod2,
				entity.OrderPaymentMethod3,
				entity.OrderPaymentMethod4,
			)),
		)
	}
	if len(i.Parcels) > 0 {
		rules = append(rules, validation.Field(&i.Parcels))
	}
	if len(i.Products) > 0 {
		rules = append(rules, validation.Field(&i.Products))
	}
	return validation.ValidateStruct(&i, rules...)
}

//Validate validate input from get order parcel API.
func (i GetOrderParcelInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&i.EndDate, validation.Required, validation.Date("2006-01-02")),
	)
}

//Validate validate input from get order product API.
func (i GetOrderProductInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&i.EndDate, validation.Required, validation.Date("2006-01-02")),
	)
}

//Validate
func (i UpdateOrderParcelInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ParcelType, validation.Required),
		validation.Field(&i.OrderParcelId, validation.Required, is.UUIDv4),
		validation.Field(&i.TrackingCode, validation.Required),
	)
}

//MakeOrder create order, order parcel, order an_parcel and order product.
func (s *service) MakeOrder(input MakeOrderInput, userId string) (MakeOrderOutput, error) {
	//Remember count of parcels and products.
	parcelCount := len(input.Parcels)
	anParcelCount := len(input.AnParcels)
	productCount := len(input.Products)

	if err := input.Validate(); err != nil {
		return MakeOrderOutput{}, internal.ErrInvalidInput{InternalError: err, Details: err}
	}

	//Initialize output.
	output := MakeOrderOutput{}
	parcelStatus := true //parcelStatus tell created all orders with no problem.
	anParcelStatus := false
	productStatus := true                                       //productStatus tell created all products with no problem.
	borList := make([]provider.BookingOrderResult, parcelCount) //BOR stand for BookingOrderResult. Have the same length with parcel input.
	anBorList := make([]provider.BookingOrderResult, anParcelCount)

	//PRODUCT, check in DATABASE.
	//Check if input have products and can make order by provider(parcelStatus). Check request products EXIST on DB.
	//Check product first because it's check on our own DB. It's faster than parcel that must HTTP request to provider.
	if productCount > 0 {
		//Allocate product output.
		output.Products = make([]makeOrderOutputStatus, productCount)
		//Split id list from product input.
		ids := make([]int32, productCount)
		for i := 0; i < productCount; i++ {
			ids[i] = input.Products[i].Id
		}
		existList, err := s.productRepository.DoesIdsExistByUserId(ids, userId)
		if err != nil {
			return MakeOrderOutput{}, err
		}
		for i, e := range existList {
			productStatus = productStatus && e
			//If doesn't exist set output message.
			if e == true {
				output.Products[i].Status = true
			} else {
				output.Products[i].Status = false
				output.Products[i].Message = fmt.Sprintf("Product id %v not found", ids[i])
			}
		}
	}

	if !productStatus {
		return output, nil
	}

	//JNA ADDRESS as ORIGIN
	//When booking and get-price use JNA ADDRESS as ORIGIN.

	//PARCEL, request order with PROVIdER.
	//Check have parcels and productStatus is true. If not, not make a request to provider for optimize time.
	if parcelCount > 0 {
		//Allocate parcel output.
		output.Parcels = make([]makeOrderOutputParcel, parcelCount)
		//Booking order. Check can booking order with PROVIDER.
		for i, parcel := range input.Parcels {
			p, err := s.providers.Get(parcel.ProviderCode)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			if s.jnaAddressAsOrigin {
				//Swap to use JNA Address when BOOKING
				parcel.Origin = s.jnaAddress
			}
			result, err := p.BookingOrder(parcel)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			borList[i] = result
			parcelStatus = parcelStatus && result.Status
		}

		//Confirm Order. If booking order is OK Confirm Order. If not cancel order.
		if parcelStatus == true {
			//CONFIRM
			//Confirm first and GetPrice in second.
			//Confirm Order. And set status and message from BookingOrderResult to output.
			for i := 0; i < parcelCount; i++ {
				p, err := s.providers.Get(borList[i].ProviderCode)
				if err != nil {
					return MakeOrderOutput{}, err
				}
				if err := p.ConfirmOrder(&borList[i]); err != nil {
					return MakeOrderOutput{}, err
				}
				//Set status and message to output.
				output.Parcels[i].Status = borList[i].Status
				output.Parcels[i].Message = borList[i].Message
				output.Parcels[i].TrackingCode = borList[i].TrackingCode
				if borList[i].ProviderCode == entity.ProviderCodeShippop {
					if borList[i].ResultShippop.CourierCode == "FLE" {
						output.Parcels[i].ShippopFlashSortingCode = &borList[i].ResultShippop.FlashSortingCode
					}
				}
			}

			//GET PRICE
			for i := 0; i < parcelCount; i++ {
				p, err := s.providers.Get(borList[i].ProviderCode)
				if err != nil {
					return MakeOrderOutput{}, err
				}
				//Get price if success populate price to output.
				//If failed combine get price result's message wit cor's message.
				if output.Parcels[i].Status == true {
					//Temporary origin
					oldOrigin := input.Parcels[i].Origin
					//Swap to use JNA Address when get price.
					if s.jnaAddressAsOrigin {
						input.Parcels[i].Origin = s.jnaAddress
					}

					r, err := p.GetPrice(input.Parcels[i])

					if s.jnaAddressAsOrigin {
						input.Parcels[i].Origin = oldOrigin
					}

					if err != nil {
						return MakeOrderOutput{}, err
					}
					//Populate price to input and output. Input for save to DB, output for user.
					if r.Status == true {
						output.Parcels[i].Status = true
						input.Parcels[i].Price = r.Price
						output.Parcels[i].Price = r.Price
					} else {
						output.Parcels[i].Status = false
						output.Parcels[i].Message += r.Message
					}
				}
			}
		} else {
			//Cancel Order.
			for i := 0; i < parcelCount; i++ {
				p, err := s.providers.Get(borList[i].ProviderCode)
				if err != nil {
					return MakeOrderOutput{}, err
				}
				if err := p.CancelOrder(borList[i]); err != nil {
					return MakeOrderOutput{}, err
				}
				output.Parcels[i].Status = borList[i].Status
				output.Parcels[i].Message = borList[i].Message
			}
		}
	}

	//AgentNetwork Parcel
	if anParcelCount > 0 {
		output.AnParcels = make([]makeOrderOutputAnParcel, anParcelCount)
		//Make group and remember order of parcel by Key
		parcels := make([]entity.Parcel, anParcelCount)
		for i := range input.AnParcels {
			parcels[i] = input.AnParcels[i].Parcel
			parcels[i].Key = i
			if s.jnaAddressAsOrigin {
				//Swap to use JNA Address when BOOKING
				parcels[i].Origin = s.jnaAddress
			}
		}
		groupParcels, err := s.providers.GroupingProvider(parcels)
		if err != nil {
			return MakeOrderOutput{}, nil
		}
		//For each provider
		for i := range groupParcels {
			p, err := s.providers.Get(i)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			//BOOKING
			results, err := p.BookingOrders(parcels)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			for _, v := range results {
				anParcelStatus = anParcelStatus || v.Status
			}
			//CONFIRM
			if anParcelStatus == true {
				if err := p.ConfirmOrders(results); err != nil {
					return MakeOrderOutput{}, err
				}
			}
			for _, v := range results {
				anBorList[v.ParcelKey] = v
			}
		}

		for i := range anBorList {
			//Put status, trackingCode and message to output.
			output.AnParcels[i].Status = anBorList[i].Status && anBorList[i].TrackingCode != ""
			output.AnParcels[i].Message = anBorList[i].Message
			output.AnParcels[i].TrackingCode = anBorList[i].TrackingCode

			if anBorList[i].ProviderCode == entity.ProviderCodeShippop {
				if anBorList[i].ResultShippop.CourierCode == "FLE" {
					output.AnParcels[i].ShippopFlashSortingCode = &anBorList[i].ResultShippop.FlashSortingCode
				}
			}
		}
	}

	//SAVE PRODUCT, PARCEL, ORIGIN and DESTINATION
	//output.Status come from two statuses are parcelStatus and productStatus.
	output.Status = parcelStatus && productStatus && anParcelStatus
	//Check if request parcels are OK and product in db are EXIST then SAVE data.
	if output.Status == true {
		var sd entity.Sender
		if parcelCount > 0 || productCount > 0 {
			//Save sender.
			var err error
			sd, err = s.senderRepository.SaveByPhoneNumber(input.Sender, input.Sender.PhoneNumber)
			if err != nil {
				return MakeOrderOutput{}, err
			}
		}

		//Save origin, destination, orderParcel and order of each provider.
		if parcelCount > 0 {
			//Split origin and destination from input parcels.
			org := make([]entity.Origin, parcelCount)
			des := make([]entity.Destination, parcelCount)
			for i := 0; i < parcelCount; i++ {
				org[i] = input.Parcels[i].Origin
				des[i] = input.Parcels[i].Destination
			}
			//Save origin, destination.
			org, err := s.originRepository.UpsertByPhoneNumber(org)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			des, err = s.destinationRepository.UpsertByPhoneNumber(des)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			//Parcel
			orderParcel := make([]entity.OrderParcel, parcelCount)
			for i := 0; i < parcelCount; i++ {
				orderParcel[i].OrderParcelId = uuid.NewV4().String()
				orderParcel[i].UserId = userId
				orderParcel[i].SenderId = sd.Id
				orderParcel[i].OriginId = org[i].Id
				orderParcel[i].DestinationId = des[i].Id
				orderParcel[i].ProviderCode = input.Parcels[i].ProviderCode
				orderParcel[i].Price = input.Parcels[i].Price
				orderParcel[i].ParcelShape = input.Parcels[i].ParcelShape
				orderParcel[i].PaymentMethod = input.PaymentMethod
				orderParcel[i].TrackingCode = output.Parcels[i].TrackingCode
				orderParcel[i].CodAmount = input.Parcels[i].CODAmount
			}
			orderParcel, err = s.orderParcelRepository.Insert(orderParcel)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			//Group each order to each provider.
			group := make(map[entity.ProviderCode][]provider.InsertOrderInput)
			for i, v := range borList {
				group[v.ProviderCode] = append(group[v.ProviderCode], provider.InsertOrderInput{
					OrderParcelId:      orderParcel[i].OrderParcelId,
					BookingOrderResult: v,
				})
			}
			//Save order of each provider eg. order_parcel_shippop.
			for k, v := range group {
				p, err := s.providers.Get(k)
				if err != nil {
					return MakeOrderOutput{}, err
				}
				if err := p.InsertOrder(v); err != nil {
					return MakeOrderOutput{}, err
				}
			}
		}

		if anParcelCount > 0 {
			//Split origin and destination from input parcels.
			org := make([]entity.Origin, anParcelCount)
			des := make([]entity.Destination, anParcelCount)
			for i := 0; i < anParcelCount; i++ {
				org[i] = input.AnParcels[i].Origin
				des[i] = input.AnParcels[i].Destination
			}
			//Save origin, destination.
			org, err := s.originRepository.UpsertByPhoneNumber(org)
			if err != nil {
				return MakeOrderOutput{}, err
			}
			des, err = s.destinationRepository.UpsertByPhoneNumber(des)
			if err != nil {
				return MakeOrderOutput{}, err
			}

			//Parcel
			var orderParcel []entity.OrderParcel
			//Group each order to each provider.
			group := make(map[entity.ProviderCode][]provider.InsertOrderInput)
			//AnParcel
			for i := 0; i < anParcelCount; i++ {
				if output.AnParcels[i].Status {
					id := uuid.NewV4().String()
					orderParcel = append(orderParcel, entity.OrderParcel{
						//Set order for save to DB.
						OrderParcelId: id,
						UserId:        userId,
						OriginId:      org[i].Id,
						DestinationId: des[i].Id,
						ProviderCode:  input.AnParcels[i].ProviderCode,
						ParcelShape:   input.AnParcels[i].ParcelShape,
						PaymentMethod: input.PaymentMethod,
						CodAmount:     input.AnParcels[i].CODAmount,
						TrackingCode:  output.AnParcels[i].TrackingCode,
					})

					//Set output
					output.AnParcels[i].OrderParcelId = id

					group[anBorList[i].ProviderCode] = append(group[anBorList[i].ProviderCode], provider.InsertOrderInput{
						OrderParcelId:      id,
						BookingOrderResult: anBorList[i],
					})
				}
			}
			orderParcel, err = s.orderParcelRepository.Insert(orderParcel)
			if err != nil {
				return MakeOrderOutput{}, err
			}

			//Save order of each provider eg. order_parcel_shippop.
			for k, v := range group {
				p, err := s.providers.Get(k)
				if err != nil {
					return MakeOrderOutput{}, err
				}
				if err := p.InsertOrder(v); err != nil {
					return MakeOrderOutput{}, err
				}
			}
		}

		//Save orderProduct.
		if productCount > 0 {
			odp := make([]entity.OrderProduct, productCount)
			//Convert to entity.OrderProduct
			for i := range odp {
				odp[i].UserId = userId
				odp[i].SenderId = sd.Id
				odp[i].ProductId = input.Products[i].Id
				odp[i].Quantity = input.Products[i].Quantity
				odp[i].PaymentMethod = input.PaymentMethod
			}
			odp, err := s.orderProductRepository.Insert(odp)
			if err != nil {
				return MakeOrderOutput{}, err
			}
		}
	}

	return output, nil
}

//GetParcelPrice check parcel price with provider and return it.
func (s *service) GetParcelPrice(input GetParcelPriceInput) (GetParcelPriceOutput, error) {
	//Validate input.
	count := len(input)
	if count == 0 {
		return GetParcelPriceOutput{}, internal.ErrInvalidInput{Details: "must have at least 1 parcel or 1 product."}
	}
	if err := validation.Validate(input); err != nil {
		return GetParcelPriceOutput{}, internal.ErrInvalidInput{InternalError: err, Details: err}
	}
	//Get Price.
	output := make(GetParcelPriceOutput, count)
	for i, parcel := range input {
		p, err := s.providers.Get(parcel.ProviderCode)
		if err != nil {
			return GetParcelPriceOutput{}, err
		}
		//Get price if success populate price to output. If failed show error message.
		r, err := p.GetPrice(input[i])
		if err != nil {
			return GetParcelPriceOutput{}, err
		}
		if r.Status == true {
			output[i].Status = true
			output[i].Price = r.Price
		} else {
			output[i].Status = false
			output[i].Message = r.Message
		}
	}

	return output, nil
}

//GetOrderParcel return order parcel.
func (s *service) GetOrderParcel(input GetOrderParcelInput) (GetOrderParcelOutput, error) {
	if err := input.Validate(); err != nil {
		return GetOrderParcelOutput{}, internal.ErrInvalidInput{InternalError: err, Details: err}
	}
	return s.orderParcelRepository.Find(input.UserId, input.StartDate, input.EndDate)
}

//GetOrderProduct return order product.
func (s *service) GetOrderProduct(input GetOrderProductInput) (GetOrderProductOutput, error) {
	if err := input.Validate(); err != nil {
		return GetOrderProductOutput{}, internal.ErrInvalidInput{InternalError: err, Details: err}
	}
	return s.orderProductRepository.Find(input.UserId, input.StartDate, input.EndDate)
}

//GetOrderParcelByIds return map of record, key=id and value=order.
func (s *service) GetOrderParcelByIds(input GetOrderParcelByIdsInput) (GetOrderParcelByIdsOutput, error) {
	err := validation.Validate(input.UserId, is.UUIDv4)
	if err != nil {
		return nil, internal.ErrInvalidInput{InternalError: err, Details: err}
	}
	err = validation.Validate(input.Ids, validation.Required, validation.Each(validation.Required), validation.Each(is.UUIDv4))
	if err != nil {
		return nil, internal.ErrInvalidInput{InternalError: err, Details: err}
	}

	return s.orderParcelRepository.FindOrderByIds(input.UserId, input.Ids)
}

//UpdateOrderParcel
func (s *service) UpdateOrderParcel(input UpdateOrderParcelInput) error {
	if err := input.Validate(); err != nil {
		return internal.ErrInvalidInput{InternalError: err, Details: err}
	}
	if err := s.orderParcelRepository.UpdateOrder(entity.OrderParcel{
		OrderParcelId: input.OrderParcelId,
		TrackingCode:  input.TrackingCode,
	}); err != nil {
		return err
	}
	//If parcel type is AgentNetwork Parcel
	if input.ParcelType == 2 {
		_ = s.SendOrderStatus(StatusHookPayload{
			SpOrderParcelId: input.OrderParcelId,
			TrackingCode:    input.TrackingCode,
		})
	}
	return nil
}

//UpdateOrderParcelStatus update order parcels status when it call by provider webhook.
func (s *service) UpdateOrderParcelStatus(input StatusData, providerCode entity.ProviderCode) error {
	//Update status in database.
	p, err := s.providers.Get(providerCode)
	if err != nil {
		return err
	}
	//In this case, input.TrackingCode is Shippop's TrackingCode not courier tracking code.
	//Update status handle about status, cod_status, status_completed_date and cod_transferred_date.
	orderParcelId, statusCompletedDate, codTransferredDate, err := p.UpdateOrderStatus(input.TrackingCode, input.Status, input.CodStatus)
	if err != nil {
		return err
	}
	//UpdateOrder handle about price, weight, width, length and height.
	if err := s.orderParcelRepository.UpdateOrder(entity.OrderParcel{
		OrderParcelId: orderParcelId,
		Price:         input.Price,
		ParcelShape: entity.ParcelShape{
			Weight: float32(input.Weight),
			Width:  float32(input.Width),
			Length: float32(input.Length),
			Height: float32(input.Height),
		},
	}); err != nil {
		return err
	}
	//Forward status to Agent-Network
	StatusHookPayload := StatusHookPayload{
		SpOrderParcelId:     orderParcelId,
		Status:              input.Status,
		CodStatus:           input.CodStatus,
		StatusCompletedDate: statusCompletedDate,
		CodTransferredDate:  codTransferredDate,
	}
	_ = s.SendOrderStatus(StatusHookPayload)
	return nil
}

//SendOrderWeight send order weight data to AgentNetwork-Service
func (s *service) SendOrderWeight(p WeightHookPayload) error {
	url := fmt.Sprintf("%v/closed/hook/order-weight", s.anUrl)
	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(p); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, &payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 1}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

//SendOrderStatus send order status data to AgentNetwork-Service
func (s *service) SendOrderStatus(p StatusHookPayload) error {
	url := fmt.Sprintf("%v/closed/hook/order-status", s.anUrl)
	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(p); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, &payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 2 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
