package courier

import (
	"encoding/json"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"os"
)

type repository struct {
	File     *os.File
	Couriers []entity.Courier
}

func NewJSONRepository(file *os.File) (*repository, error) {
	var courier []entity.Courier
	if err := json.NewDecoder(file).Decode(&courier); err != nil {
		return nil, err
	}
	return &repository{
		Couriers: courier,
	}, nil
}

func (r *repository) FindAll() ([]entity.Courier, error) {
	return r.Couriers, nil
}
