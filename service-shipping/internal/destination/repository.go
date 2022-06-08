package destination

import "github.com/jna-distribution/service-shipping/internal/entity"

type Repository interface {
	DoesPhoneNumberExist(phoneNumber string) (bool, error)
	UpsertByPhoneNumber(destinations []entity.Destination) ([]entity.Destination, error)
}
