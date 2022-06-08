package handler

import (
	"github.com/jna-distribution/service-shipping/cmd/restapi/authentication"
	"github.com/jna-distribution/service-shipping/cmd/restapi/middleware"
	"github.com/jna-distribution/service-shipping/cmd/restapi/responses"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/courier"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/product"
	"net/http"
)

func BuildCourierProductHandler(rootMux *http.ServeMux, courierService courier.Service, productService product.Service) {
	rootMux.Handle("/courier-and-product/", middleware.WithPermission(courierAndProductHandler(courierService, productService), entity.RoleShop, entity.RoleAgentNetworkMember))
}

type courierAndProductResponse struct {
	Courier []entity.Courier `json:"courier"`
	Product []entity.Product `json:"product"`
}

func courierAndProductHandler(courierService courier.Service, productService product.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		if r.RequestURI != "/courier-and-product/" {
			e := responses.BuildErrorResponse(internal.ErrNotFound{})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		if r.Method == http.MethodGet {
			claims, err := authentication.DetachClaimsFromHeader(r)
			c, err := courierService.Get()
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			p, err := productService.GetByContactId(claims.UserId)
			if err != nil {
				e := responses.BuildErrorResponse(err)
				rw.WriteWithStatus(e.Status, e)
				return
			}
			res := courierAndProductResponse{Courier: c, Product: p}
			rw.WriteWithStatus(http.StatusOK, res)
			return
		} else {
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
			return
		}
	})
}
