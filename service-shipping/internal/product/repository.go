package product

import "github.com/jna-distribution/service-shipping/internal/entity"

type Repository interface {
	Insert(product entity.Product) (entity.Product, error)
	DoesIdsExistByUserId(ids []int32, userId string) ([]bool, error)
	FindAll() ([]entity.Product, error)
	FindByContactId(userId string) ([]entity.Product, error)
	FindByContactIdAndId(userId string, id int32) (entity.Product, error)
	Update(product entity.Product) (entity.Product, error)
	Delete(product entity.Product) error
}
