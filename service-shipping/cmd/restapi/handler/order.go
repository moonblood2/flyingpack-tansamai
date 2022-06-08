package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jna-distribution/service-shipping/cmd/restapi/authentication"
	"github.com/jna-distribution/service-shipping/cmd/restapi/middleware"
	"github.com/jna-distribution/service-shipping/cmd/restapi/requests"
	"github.com/jna-distribution/service-shipping/cmd/restapi/responses"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/order"
)

const (
	parcelPricePattern = "/parcel-price/"

	orderPattern                        = "/orders/"
	orderParcelPattern                  = "/orders/parcels/"
	orderParcelByDatePattern            = "/orders/parcels/start/:start_date/end/:end_date/"
	orderParcelByIds                    = "/orders/parcels/ids/"
	orderProductPattern                 = "/orders/products/"
	orderProductByDatePattern           = "/orders/products/start/:start_date/end/:end_date/"
	orderStatusHookPattern              = "/order/status/hook/"
	orderStatusHookByCourierNamePattern = "/order/status/hook/:provider"
)

//phe: Provider hook endpoint
const (
	pheShippop = "rIDJUMrAtESIr9Ge8uB6"
)

func BuildOrderHandler(rootMux *http.ServeMux, orderSvc order.Service) {
	rootMux.Handle(parcelPricePattern, middleware.WithPermission(getParcelPriceHandler(orderSvc), entity.RoleShop, entity.RoleAgentNetworkMember))
	//Order mix
	rootMux.Handle(orderPattern, middleware.WithPermission(orderHandler(orderSvc), entity.RoleShop, entity.RoleAgentNetworkMember))
	//Order parcel
	rootMux.Handle(orderParcelPattern, middleware.WithPermission(orderParcelHandler(orderSvc), entity.RoleShop, entity.RoleAgentNetworkMember))
	rootMux.Handle(orderParcelByIds, middleware.WithPermission(orderParcelByIdsHandler(orderSvc), entity.RoleShop, entity.RoleAgentNetworkMember, entity.RoleAccounting))
	//Order product
	rootMux.Handle(orderProductPattern, middleware.WithPermission(orderProductHandler(orderSvc), entity.RoleShop))

	//Public, web hook
	rootMux.Handle(orderStatusHookPattern, orderStatusHook(orderSvc))
}

func getParcelPriceHandler(orderSvc order.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		if r.RequestURI != parcelPricePattern {
			e := responses.BuildErrorResponse(internal.ErrNotFound{})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		switch r.Method {
		case http.MethodPost:
			//Retrieve data from JSON body.
			input := order.GetParcelPriceInput{}
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			output, err := orderSvc.GetParcelPrice(input)
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			rw.WriteWithStatus(http.StatusOK, output)
		default:
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
		}
	})
}

func orderHandler(orderSvc order.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		if r.RequestURI != orderPattern {
			e := responses.BuildErrorResponse(internal.ErrNotFound{})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		switch r.Method {
		case http.MethodPost:
			createOrder(w, r, orderSvc)
		default:
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
		}
	})
}

func orderParcelHandler(orderSvc order.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		switch r.Method {
		case http.MethodGet:
			getOrderParcel(w, r, orderSvc)
		case http.MethodPut:
			updateOrderParcel(w, r, orderSvc)
		default:
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
		}
	})
}

func orderProductHandler(orderSvc order.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		switch r.Method {
		case http.MethodGet:
			getOrderProduct(w, r, orderSvc)
		default:
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
		}
	})
}

func orderParcelByIdsHandler(orderSvc order.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		switch r.Method {
		case http.MethodPost:
			getOrderParcelByIds(w, r, orderSvc)
		default:
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
		}
	})
}

func createOrder(w http.ResponseWriter, r *http.Request, orderSvc order.Service) {
	rw := responses.NewJSONResponseWriter(w)
	//Retrieve data from JSON body.
	input := order.MakeOrderInput{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	//Retrieve user id from claim.
	claims, err := authentication.DetachClaimsFromHeader(r)
	if err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	output, err := orderSvc.MakeOrder(input, claims.UserId)
	if err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	rw.WriteWithStatus(http.StatusOK, output)
}

func getOrderParcel(w http.ResponseWriter, r *http.Request, orderSvc order.Service) {
	rw := responses.NewJSONResponseWriter(w)
	if requests.MatchPattern(orderParcelByDatePattern, r.RequestURI) {
		param, err := requests.NewPathParam(orderParcelByDatePattern, r.RequestURI)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		claims, err := authentication.DetachClaimsFromHeader(r)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		input := order.GetOrderParcelInput{
			UserId:    claims.UserId,
			StartDate: param["start_date"],
			EndDate:   param["end_date"],
		}
		output, err := orderSvc.GetOrderParcel(input)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		rw.WriteWithStatus(http.StatusOK, output)
	} else {
		e := responses.BuildErrorResponse(internal.ErrNotFound{})
		rw.WriteWithStatus(e.Status, e)
		return
	}
}

func getOrderParcelByIds(w http.ResponseWriter, r *http.Request, orderSvc order.Service) {
	rw := responses.NewJSONResponseWriter(w)
	var input order.GetOrderParcelByIdsInput
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	claims, err := authentication.DetachClaimsFromHeader(r)
	if err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	if input.UserId == "" {
		input.UserId = claims.UserId
	}
	orders, err := orderSvc.GetOrderParcelByIds(input)
	if err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	rw.WriteWithStatus(http.StatusOK, orders)
}

func updateOrderParcel(w http.ResponseWriter, r *http.Request, orderSvc order.Service) {
	rw := responses.NewJSONResponseWriter(w)
	var input order.UpdateOrderParcelInput
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	err := orderSvc.UpdateOrderParcel(input)
	if err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	rw.WriteWithStatus(http.StatusOK, struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{true, "success"})
}

func getOrderProduct(w http.ResponseWriter, r *http.Request, orderSvc order.Service) {
	rw := responses.NewJSONResponseWriter(w)
	if requests.MatchPattern(orderProductByDatePattern, r.RequestURI) {
		param, err := requests.NewPathParam(orderProductByDatePattern, r.RequestURI)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		claims, err := authentication.DetachClaimsFromHeader(r)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		input := order.GetOrderProductInput{
			UserId:    claims.UserId,
			StartDate: param["start_date"],
			EndDate:   param["end_date"],
		}
		output, err := orderSvc.GetOrderProduct(input)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		rw.WriteWithStatus(http.StatusOK, output)
	} else {
		e := responses.BuildErrorResponse(internal.ErrNotFound{})
		rw.WriteWithStatus(e.Status, e)
		return
	}
}

func orderStatusHook(orderSvc order.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		if requests.MatchPattern(orderStatusHookByCourierNamePattern, r.RequestURI) {
			param, err := requests.NewPathParam(orderStatusHookByCourierNamePattern, r.RequestURI)
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			pro := param["provider"]
			orderStatus := order.StatusData{}
			switch pro {
			case pheShippop:
				defer r.Body.Close()
				//headerContentTtype := r.Header.Get("Content-Type")
				//log.Println(headerContentTtype)
				r.ParseForm()
				orderStatus.TrackingCode = r.PostForm.Get("tracking_code")
				status := r.PostForm.Get("order_status")
				if status == "wait" || status == "booking" || status == "shipping" || status == "complete" || status == "cancel" || status == "return" {
					orderStatus.Status = status
				} else if status == "pending_transfer" || status == "transferred" || status == "cancel_transfer" {
					orderStatus.CodStatus = status
				}
				if status == "shipping" {
					weightStr := r.PostForm.Get("data[weight]")
					widthStr := r.PostForm.Get("data[width]")
					heightStr := r.PostForm.Get("data[height]")
					lengthStr := r.PostForm.Get("data[length]")
					priceStr := r.PostForm.Get("data[price]")
					dateTimeStr := r.PostForm.Get("data[datetime]")

					weight, _ := strconv.ParseFloat(weightStr, 32)
					width, _ := strconv.ParseFloat(widthStr, 32)
					height, _ := strconv.ParseFloat(heightStr, 32)
					length, _ := strconv.ParseFloat(lengthStr, 32)
					price, _ := strconv.ParseFloat(priceStr, 32)
					dateTime, _ := time.Parse("2006-01-02 15:04:05", dateTimeStr)

					orderStatus.Weight = weight
					orderStatus.Width = width
					orderStatus.Height = height
					orderStatus.Length = length
					orderStatus.Price = price
					orderStatus.Datetime = dateTime
				}
				//If not have both status and codStatus.
				if orderStatus.Status == "" && orderStatus.CodStatus == "" {
					rw.WriteWithStatus(http.StatusBadRequest, struct {
						Status bool `json:"status"`
					}{false})
					return
				}

				err := orderSvc.UpdateOrderParcelStatus(orderStatus, entity.ProviderCodeShippop)
				if err != nil {
					e := responses.BuildErrorResponse(err)
					rw.WriteWithStatus(e.Status, e)
					return
				}
				rw.WriteWithStatus(http.StatusOK, struct {
					Status bool `json:"status"`
				}{true})
			default:
				log.Printf("phe: %v not match", pro)
				rw.WriteWithStatus(http.StatusBadRequest, struct {
					Status bool `json:"status"`
				}{false})
			}
		}
	})
}
