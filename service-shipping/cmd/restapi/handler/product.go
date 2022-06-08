package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jna-distribution/service-shipping/cmd/restapi/authentication"
	"github.com/jna-distribution/service-shipping/cmd/restapi/middleware"
	"github.com/jna-distribution/service-shipping/cmd/restapi/requests"
	"github.com/jna-distribution/service-shipping/cmd/restapi/responses"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/product"
	"net/http"
	"strconv"
)

const productPattern = "/product/"
const productByIdPattern = "/product/:id/"

func BuildProductHandler(rootMux *http.ServeMux, service product.Service) {
	// /product/ -> Get
	// /product/:id/ -> Get by id, Edit by id, Del by id.
	rootMux.Handle(productPattern, middleware.WithPermission(productHandler(service), entity.RoleShop))
}

func productHandler(service product.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		switch r.Method {
		case http.MethodPost:
			createProduct(w, r, service)
		case http.MethodGet:
			getProduct(w, r, service)
		case http.MethodPut:
			editProduct(w, r, service)
		case http.MethodDelete:
			removeProduct(w, r, service)
		default:
			e := responses.BuildErrorResponse(internal.ErrNotFound{})
			rw.WriteWithStatus(e.Status, e)
			return
		}
	})
}

type createProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type createProductResponse struct {
	entity.Product
}

func createProduct(w http.ResponseWriter, r *http.Request, service product.Service) {
	rw := responses.NewJSONResponseWriter(w)
	if r.RequestURI != productPattern {
		e := responses.BuildErrorResponse(internal.ErrNotFound{})
		rw.WriteWithStatus(e.Status, e)
		return
	}
	req := createProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := responses.BuildErrorResponse(responses.ErrBadRequest{})
		rw.WriteWithStatus(e.Status, e)
		return
	}
	claims, err := authentication.DetachClaimsFromHeader(r)
	if err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	p, err := service.Create(entity.Product{
		UserId: claims.UserId,
		Name:   req.Name,
		Price:  req.Price,
	})
	if err != nil {
		e := responses.BuildErrorResponse(err)
		rw.WriteWithStatus(e.Status, e)
		return
	}
	rw.WriteWithStatus(http.StatusOK, createProductResponse{p})
}

type getProductResponse struct {
	Product []entity.Product `json:"product"`
}

type getProductByIdResponse struct {
	Product *entity.Product `json:"product"`
}

func getProduct(w http.ResponseWriter, r *http.Request, service product.Service) {
	rw := responses.NewJSONResponseWriter(w)
	//Get All
	if r.RequestURI == productPattern {
		claims, err := authentication.DetachClaimsFromHeader(r)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		p, err := service.GetByContactId(claims.UserId)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		rw.WriteWithStatus(http.StatusOK, getProductResponse{p})
	} else if requests.MatchPattern(productByIdPattern, r.RequestURI) {
		param, err := requests.NewPathParam(productByIdPattern, r.RequestURI)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		idStr := param["id"]
		id, err := strconv.ParseInt(idStr, 10, 32)
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
		p, err := service.GetByContactIdAndId(claims.UserId, int32(id))
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		rw.WriteWithStatus(http.StatusOK, getProductByIdResponse{&p})
		return
	} else {
		e := responses.BuildErrorResponse(internal.ErrNotFound{})
		rw.WriteWithStatus(e.Status, e)
		return
	}
}

type editProductByIdRequest struct {
	entity.Product
}

type editProductByIdResponse struct {
	entity.Product
}

func editProduct(w http.ResponseWriter, r *http.Request, service product.Service) {
	rw := responses.NewJSONResponseWriter(w)
	if requests.MatchPattern(productByIdPattern, r.RequestURI) {
		param, err := requests.NewPathParam(productByIdPattern, r.RequestURI)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		idStr := param["id"]
		id, err := strconv.ParseInt(idStr, 10, 32)
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
		req := editProductByIdRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			e := responses.BuildErrorResponse(responses.ErrBadRequest{InternalError: err})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		p, err := service.Edit(entity.Product{
			Id:     int32(id),
			UserId: claims.UserId,
			Name:   req.Name,
			Price:  req.Price,
		})
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		rw.WriteWithStatus(http.StatusOK, editProductByIdResponse{p})
	} else {
		e := responses.BuildErrorResponse(internal.ErrNotFound{})
		rw.WriteWithStatus(e.Status, e)
		return
	}
}

type removeProductByIdResponse struct {
	Message string `json:"message"`
}

func removeProduct(w http.ResponseWriter, r *http.Request, service product.Service) {
	rw := responses.NewJSONResponseWriter(w)
	if requests.MatchPattern(productByIdPattern, r.RequestURI) {
		param, err := requests.NewPathParam(productByIdPattern, r.RequestURI)
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		idStr := param["id"]
		id, err := strconv.ParseInt(idStr, 10, 32)
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
		err = service.Remove(entity.Product{
			Id:     int32(id),
			UserId: claims.UserId,
		})
		if err != nil {
			e := responses.BuildErrorResponse(err)
			rw.WriteWithStatus(e.Status, e)
			return
		}
		rw.WriteWithStatus(http.StatusOK, removeProductByIdResponse{fmt.Sprintf("Product %v has deleted.", id)})
	} else {
		e := responses.BuildErrorResponse(internal.ErrNotFound{})
		rw.WriteWithStatus(e.Status, e)
		return
	}
}
