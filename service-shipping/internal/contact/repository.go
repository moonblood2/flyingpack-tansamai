package contact

import "github.com/jna-distribution/service-shipping/internal/entity"

type Repository interface {
	DoesPhoneNumberExist(phoneNumber string) (bool, error)
	FindByUserId(userId string) (entity.Contact, error)
	Create(contact entity.Contact) (entity.Contact, error)
	CreateWithUser(user entity.User, contact entity.Contact) (entity.User, entity.Contact, error)
}
