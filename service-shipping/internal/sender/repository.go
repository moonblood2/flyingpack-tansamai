package sender

import "github.com/jna-distribution/service-shipping/internal/entity"

type Repository interface {
	DoesIdExist(id int64) (bool, error)
	DoesPhoneNumberExist(phoneNumber string) (bool, error)
	FindById(id int64) (entity.Sender, error)
	FindByPhoneNumber(phoneNumber string) (entity.Sender, error)
	Create(sender entity.Sender) (entity.Sender, error)
	UpdateById(sender entity.Sender, id int64) (entity.Sender, error)
	UpdateByPhoneNumber(sender entity.Sender, phoneNumber string) (entity.Sender, error)
	//SaveByPhoneNumber create when phoneNumber doesn't exists and update in the other hand.
	SaveByPhoneNumber(sender entity.Sender, phoneNumber string) (entity.Sender, error)
}
