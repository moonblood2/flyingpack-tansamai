package courier

import "github.com/jna-distribution/service-shipping/internal/entity"

type Repository interface {
	FindAll() ([]entity.Courier, error)
}
