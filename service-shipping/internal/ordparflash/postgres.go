package ordparflash

import (
	"database/sql"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"strings"
	"time"
)

type repository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) Insert(orders []entity.OrderParcelFlash) ([]entity.OrderParcelFlash, error) {
	//Make value placeholder list and value list.
	vpl := make([]string, len(orders))
	var values []interface{}
	//n for count number of placeholder($N).
	n := 0
	for i, v := range orders {
		vpl[i] = fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v)", n+1, n+2, n+3, n+4, n+5, n+6)
		n += 6

		values = append(values, v.OrderParcelId)
		values = append(values, v.Pno)
		values = append(values, v.State)
		values = append(values, v.StateText)
		values = append(values, v.CODAmount)
		values = append(values, time.Now())
	}
	query := fmt.Sprintf(`INSERT INTO public.order_parcel_flash(order_parcel_id, pno, state, state_text, cod_amount, created_at) VALUES %v RETURNING id;`, strings.Join(vpl, ", "))
	rows, err := r.DB.Query(query, values...)
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err := rows.Scan(&orders[i].Id)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		i++
	}
	if rows.Err() != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	return orders, nil
}
