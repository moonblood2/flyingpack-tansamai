package courier

import (
	"github.com/jna-distribution/service-shipping/internal/entity"
)

type Service interface {
	Get() ([]entity.Courier, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) Get() ([]entity.Courier, error) {
	return s.repository.FindAll()
}
