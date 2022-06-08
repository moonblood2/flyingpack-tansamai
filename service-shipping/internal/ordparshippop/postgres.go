package ordparshippop

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

func (r *repository) Insert(orders []entity.OrderParcelShippop) ([]entity.OrderParcelShippop, error) {
	//Insert order_parcel_shippop
	//Make value placeholder list and value list.
	vpl := make([]string, len(orders))
	var values []interface{}
	//n for count number of placeholder($N).
	n := 0
	for i, v := range orders {
		vpl[i] = fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)", n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8)
		n += 8

		values = append(values, v.OrderParcelId)
		values = append(values, v.PurchaseId)
		values = append(values, v.Status)
		values = append(values, v.CourierCode)
		values = append(values, v.CourierTrackingCode)
		values = append(values, v.TrackingCode)
		values = append(values, v.CODAmount)
		values = append(values, time.Now())
	}

	//Begin transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`INSERT INTO public.order_parcel_shippop(order_parcel_id, purchase_id, status, courier_code, courier_tracking_code, tracking_code, cod_amount, created_at) VALUES %v RETURNING id;`, strings.Join(vpl, ", "))
	rows, err := tx.Query(query, values...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()

	var orderParcelShippopFlashList []entity.OrderParcelShippopFlash
	var mapOrderIndex []int

	i := 0
	for rows.Next() {
		err := rows.Scan(&orders[i].Id)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, internal.ErrDatabase{InternalError: err}
			}
			return nil, internal.ErrDatabase{InternalError: err}
		}
		if orders[i].CourierCode == "FLE" {
			orders[i].OrderParcelShippopFlash.OrderParcelShippopId = orders[i].Id
			orderParcelShippopFlashList = append(orderParcelShippopFlashList, orders[i].OrderParcelShippopFlash)
			mapOrderIndex = append(mapOrderIndex, i)
		}
		i++
	}
	if rows.Err() != nil {
		if err := tx.Rollback(); err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}

	if orderParcelShippopFlashList != nil {
		//Insert order_parcel_shippop_flash
		//Make value placeholder list and value list.
		vpl := make([]string, len(orderParcelShippopFlashList))
		var values []interface{}
		//n for count number of placeholder($N).
		n := 0
		for i, v := range orderParcelShippopFlashList {
			vpl[i] = fmt.Sprintf("($%v, $%v, $%v, $%v)", n+1, n+2, n+3, n+4)
			n += 4

			values = append(values, v.OrderParcelShippopId)
			values = append(values, v.SortCode)
			values = append(values, v.DstCode)
			values = append(values, v.SortingLineCode)
		}

		query := fmt.Sprintf(`
INSERT INTO public.order_parcel_shippop_flash(
	order_parcel_shippop_id, sort_code, dst_code, sorting_line_code
) 
VALUES %v RETURNING id;`, strings.Join(vpl, ", "))
		rows, err := tx.Query(query, values...)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, internal.ErrDatabase{InternalError: err}
			}
			return nil, internal.ErrDatabase{InternalError: err}
		}
		defer rows.Close()

		i := 0
		for rows.Next() {
			err := rows.Scan(&orderParcelShippopFlashList[i].Id)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return nil, internal.ErrDatabase{InternalError: err}
				}
				return nil, internal.ErrDatabase{InternalError: err}
			}
			i++
		}
		if rows.Err() != nil {
			if err := tx.Rollback(); err != nil {
				return nil, internal.ErrDatabase{InternalError: err}
			}
			return nil, internal.ErrDatabase{InternalError: err}
		}
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}

	//Populate orderParcelShippopFlash to orders
	if orderParcelShippopFlashList != nil && mapOrderIndex != nil {
		for i, v := range orderParcelShippopFlashList {
			if orders[mapOrderIndex[i]].CourierCode == "FLE" {
				orders[i].OrderParcelShippopFlash = v
			}
		}
	}

	return orders, nil
}

func (r *repository) UpdateOrderStatus(trackingCode, status, codStatus string) (orderParcelId string, statusCompletedDate, codTransferredDate *time.Time, err error) {
	var courierTrackingCode string
	//courier_tracking_code in order_parcel_shippop = tracking_code in order_parcel
	//Find courier_tracking_code
	tx, err := r.DB.Begin()
	if err != nil {
		return orderParcelId, statusCompletedDate, codTransferredDate, internal.ErrDatabase{InternalError: err}
	}
	row := tx.QueryRow(`SELECT order_parcel_id, courier_tracking_code 
FROM public.order_parcel_shippop 
WHERE tracking_code=$1`,
		trackingCode,
	)
	err = row.Scan(&orderParcelId, &courierTrackingCode)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return orderParcelId, statusCompletedDate, codTransferredDate, internal.ErrDatabase{InternalError: err}
		}
		return orderParcelId, statusCompletedDate, codTransferredDate, internal.ErrDatabase{InternalError: err, Details: "internal:ordparshippop:postgresq:UpdateOrderStatus"}
	}
	//$1=order_parcel_id
	values := []interface{}{orderParcelId}
	var plc []string
	if status != "" {
		values = append(values, status)
		plc = append(plc, fmt.Sprintf("status=$%d", len(values)))
		if status == "complete" {
			now := time.Now()
			statusCompletedDate = &now
			values = append(values, statusCompletedDate)
			plc = append(plc, fmt.Sprintf("status_completed_date=$%d", len(values)))
		}
	}
	if codStatus != "" {
		values = append(values, codStatus)
		plc = append(plc, fmt.Sprintf("cod_status=$%d", len(values)))
		if codStatus == "transferred" {
			now := time.Now()
			codTransferredDate = &now
			values = append(values, codTransferredDate)
			plc = append(plc, fmt.Sprintf("cod_transferred_date=$%d", len(values)))
		}
	}
	if len(values) > 1 {
		//Update order_parcel.status
		query := fmt.Sprintf(`
UPDATE public.order_parcel 
SET %v
WHERE order_parcel_id=$1
`, strings.Join(plc, ","))
		_, err = tx.Exec(query, values...)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return orderParcelId, statusCompletedDate, codTransferredDate, internal.ErrDatabase{InternalError: err}
			}
			return orderParcelId, statusCompletedDate, codTransferredDate, internal.ErrDatabase{InternalError: err, Details: "internal:ordparshippop:postgresq:UpdateOrderStatus"}
		}
	}
	if err := tx.Commit(); err != nil {
		return orderParcelId, statusCompletedDate, codTransferredDate, internal.ErrDatabase{InternalError: err}
	}
	return orderParcelId, statusCompletedDate, codTransferredDate, nil
}
