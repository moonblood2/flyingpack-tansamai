package origin

import "github.com/jna-distribution/service-shipping/internal/entity"

type Repository interface {
	DoesPhoneNumberExist(phoneNumber string) (bool, error)
	DoesPhoneNumbersExist(phoneNumbers []string) ([]bool, error)
	UpsertByPhoneNumber(origins []entity.Origin) ([]entity.Origin, error)
}
