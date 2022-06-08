package ordproduct

import (
	"database/sql"
	"errors"
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

func (r *repository) Insert(orders []entity.OrderProduct) ([]entity.OrderProduct, error) {
	//Make value placeholder list and value list.
	vpl := make([]string, len(orders))
	var values []interface{}
	//n for count number of placeholder($N).
	n := 0
	for i, v := range orders {
		vpl[i] = fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v)", n+1, n+2, n+3, n+4, n+5, n+6)
		n += 6

		values = append(values, v.UserId)
		values = append(values, v.SenderId)
		values = append(values, v.ProductId)
		values = append(values, v.Quantity)
		values = append(values, v.PaymentMethod)
		values = append(values, time.Now())
	}
	//Make query.
	query := fmt.Sprintf(`INSERT INTO public.order_product (user_id, sender_id, product_id, quantity, payment_method, created_at) VALUES %v RETURNING id;`, strings.Join(vpl, ", "))
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

func (r *repository) Find(userId string, startDate, endDate string) ([]FindRecord, error) {
	rows, err := r.DB.Query(`SELECT order_product.*, sender.*, product.*
FROM order_product
INNER JOIN sender ON order_product.sender_id = sender.id AND sender.deleted_at=0
INNER JOIN product ON order_product.product_id = product.id AND product.deleted_at=0
WHERE order_product.user_id=$1 AND CAST(order_product.created_at AS DATE) BETWEEN $2 AND $3 AND order_product.deleted_at=0`,
		userId, startDate, endDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	var records []FindRecord
	for rows.Next() {
		r := FindRecord{}
		err := rows.Scan(&r.OrderProduct.Id, &r.OrderProduct.UserId, &r.OrderProduct.SenderId, &r.OrderProduct.ProductId,
			&r.OrderProduct.Quantity, &r.OrderProduct.PaymentMethod, &r.OrderProduct.CreatedAt, &r.OrderProduct.DeletedAt,
			&r.Sender.Id, &r.Sender.SenderType, &r.Sender.NationalIdNumber, &r.Sender.TaxIdNumber, &r.Sender.PassportNumber, &r.Sender.BirthDate, &r.Sender.Name, &r.Sender.PhoneNumber,
			&r.Sender.Address, &r.Sender.District, &r.Sender.State, &r.Sender.Province, &r.Sender.Postcode, &r.Sender.CreatedAt, &r.Sender.DeletedAt,
			&r.Product.Id, &r.Product.UserId, &r.Product.Name, &r.Product.Price, &r.Product.CreatedAt, &r.Product.DeletedAt,
		)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		records = append(records, r)
	}
	return records, nil
}
