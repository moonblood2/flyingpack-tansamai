package product

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
)

type Service interface {
	Create(product entity.Product) (entity.Product, error)
	GetAll() ([]entity.Product, error)
	GetByContactId(userId string) ([]entity.Product, error)
	GetByContactIdAndId(userId string, id int32) (entity.Product, error)
	Edit(product entity.Product) (entity.Product, error)
	Remove(product entity.Product) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(product entity.Product) (entity.Product, error) {
	if err := product.Validate(); err != nil {
		return entity.Product{}, internal.ErrInvalidInput{InternalError: err, Details: err}
	}
	return s.repository.Insert(product)
}

func (s *service) GetAll() ([]entity.Product, error) {
	return s.repository.FindAll()
}

func (s *service) GetByContactId(userId string) ([]entity.Product, error) {
	return s.repository.FindByContactId(userId)
}

func (s *service) GetByContactIdAndId(userId string, id int32) (entity.Product, error) {
	return s.repository.FindByContactIdAndId(userId, id)
}

func (s *service) Edit(product entity.Product) (entity.Product, error) {
	err := validation.ValidateStruct(validation.Field(&product.UserId, validation.Required, validation.Min(1)))
	if err != nil {
		return entity.Product{}, internal.ErrInvalidInput{InternalError: err, Details: nil}
	}
	err = validation.ValidateStruct(&product,
		validation.Field(&product.Id, validation.Required, validation.Min(1)),
		validation.Field(&product.Name, append(entity.ProductNameRule, validation.Required)...),
		validation.Field(&product.Price, append(entity.ProductPriceRule, validation.Required)...),
	)
	if err != nil {
		return entity.Product{}, internal.ErrInvalidInput{InternalError: err, Details: err}
	}
	return s.repository.Update(product)
}

func (s *service) Remove(product entity.Product) error {
	return s.repository.Delete(product)
}
