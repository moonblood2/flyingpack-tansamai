package user

import "github.com/jna-distribution/service-shipping/internal/entity"

// Repository interface
type Repository interface {
	DoesEmailExist(email string) (bool, error)
	Create(user entity.User) (entity.User, error)
	FindAll() ([]entity.User, error)
	FindById(id string) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	Delete(id string) error
}
