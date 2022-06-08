package ordparcel

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg/str"
	"github.com/satori/go.uuid"
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

func (r *repository) Insert(orders []entity.OrderParcel) ([]entity.OrderParcel, error) {
	if len(orders) < 1 {
		return orders, nil
	}
	//Make value placeholder list and value list.
	vpl := make([]string, len(orders))
	var values []interface{}
	//n for count number of placeholder($N).
	n := 0
	for i, v := range orders {
		vpl[i] = fmt.Sprintf("($%v::uuid, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)", n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12, n+13, n+14)
		n += 14
		if v.OrderParcelId == "" {
			v.OrderParcelId = uuid.NewV4().String()
		}
		values = append(values, v.OrderParcelId)
		values = append(values, v.UserId)
		if v.SenderId != 0 {
			values = append(values, v.SenderId)
		} else {
			values = append(values, nil)
		}
		values = append(values, v.OriginId)
		values = append(values, v.DestinationId)
		values = append(values, v.ProviderCode)
		values = append(values, v.Price)
		values = append(values, v.PaymentMethod)
		values = append(values, v.Weight)
		values = append(values, v.Width)
		values = append(values, v.Length)
		values = append(values, v.Height)
		values = append(values, v.TrackingCode)
		values = append(values, v.CodAmount)
	}
	query := fmt.Sprintf(`INSERT INTO public.order_parcel(order_parcel_id, user_id, sender_id, origin_id, destination_id, provider_code, price, payment_method, weight, width, length, height, tracking_code, cod_amount) VALUES %v;`, strings.Join(vpl, ", "))
	rows, err := r.DB.Query(query, values...)
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	return orders, nil
}

func (r *repository) UpdateOrder(order entity.OrderParcel) error {
	values := []interface{}{order.OrderParcelId}
	var setFields []string
	if order.TrackingCode != "" {
		values = append(values, order.TrackingCode)
		setFields = append(setFields, fmt.Sprintf("tracking_code=$%v", len(values)))
	}
	if order.Price != 0 {
		values = append(values, order.Price)
		setFields = append(setFields, fmt.Sprintf("price=$%v", len(values)))
	}
	if order.Weight != 0 {
		values = append(values, order.Weight)
		setFields = append(setFields, fmt.Sprintf("weight=$%v", len(values)))
	}
	if order.Width != 0 {
		values = append(values, order.Width)
		setFields = append(setFields, fmt.Sprintf("width=$%v", len(values)))
	}
	if order.Length != 0 {
		values = append(values, order.Length)
		setFields = append(setFields, fmt.Sprintf("length=$%v", len(values)))
	}
	if order.Height != 0 {
		values = append(values, order.Height)
		setFields = append(setFields, fmt.Sprintf("height=$%v", len(values)))
	}
	if len(values) > 1 {
		query := fmt.Sprintf(`UPDATE public.order_parcel SET %v WHERE order_parcel_id=$1`, strings.Join(setFields, ","))
		_, err := r.DB.Exec(query, values...)
		if err != nil {
			return internal.ErrDatabase{InternalError: err, Details: "orderparcel:postgres:UpdateOrder"}
		}
	}
	return nil
}

func (r *repository) Find(userId string, startDate, endDate string) ([]FindRecord, error) {
	rows, err := r.DB.Query(`
SELECT order_parcel.order_parcel_id, order_parcel.provider_code, order_parcel.price, order_parcel.payment_method, order_parcel.tracking_code, order_parcel.status, 
order_parcel.weight, order_parcel.width, order_parcel.length, order_parcel.height, order_parcel.cod_amount, order_parcel.created_at,
sender.*, 
origin.*, 
destination.*
FROM order_parcel
INNER JOIN sender ON order_parcel.sender_id = sender.id AND sender.deleted_at=0
INNER JOIN origin ON order_parcel.origin_id = origin.id AND origin.deleted_at=0
INNER JOIN destination ON order_parcel.destination_id = destination.id AND destination.deleted_at=0
WHERE order_parcel.user_id=$1 AND CAST(order_parcel.created_at AS DATE) BETWEEN $2 AND $3 AND order_parcel.deleted_at=0
ORDER BY order_parcel.created_at DESC`,
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
		err := rows.Scan(&r.OrderParcel.OrderParcelId, &r.OrderParcel.ProviderCode, &r.OrderParcel.Price, &r.OrderParcel.PaymentMethod, &r.OrderParcel.TrackingCode, &r.OrderParcel.Status,
			&r.OrderParcel.Weight, &r.OrderParcel.Width, &r.OrderParcel.Length, &r.OrderParcel.Height, &r.OrderParcel.CodAmount, &r.OrderParcel.CreatedAt,
			&r.Sender.Id, &r.Sender.SenderType, &r.Sender.NationalIdNumber, &r.Sender.TaxIdNumber, &r.Sender.PassportNumber, &r.Sender.BirthDate,
			&r.Sender.Name, &r.Sender.PhoneNumber, &r.Sender.Address, &r.Sender.District, &r.Sender.State, &r.Sender.Province, &r.Sender.Postcode, &r.Sender.CreatedAt, &r.Sender.DeletedAt,
			&r.Origin.Id, &r.Origin.Name, &r.Origin.PhoneNumber, &r.Origin.Address, &r.Origin.District, &r.Origin.State, &r.Origin.Province, &r.Origin.Postcode, &r.Origin.CreatedAt, &r.Origin.DeletedAt,
			&r.Destination.Id, &r.Destination.Name, &r.Destination.PhoneNumber, &r.Destination.Address, &r.Destination.District, &r.Destination.State, &r.Destination.Province, &r.Destination.Postcode, &r.Destination.CreatedAt, &r.Destination.DeletedAt,
		)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		records = append(records, r)
	}
	err = rows.Err()
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	return records, nil
}

func (r *repository) UpdateStatus(trackingCode, status string, providerCode entity.ProviderCode) error {
	_, err := r.DB.Exec(`UPDATE public.order_parcel SET status=$1 WHERE provider_code=$2 AND tracking_code=$3;`,
		status, providerCode, trackingCode,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return internal.ErrDatabase{InternalError: err}
	}
	return nil
}

func (r *repository) FindUser(trackingCode string, providerCode entity.ProviderCode) (string, error) {
	var userId string
	row := r.DB.QueryRow(
		`SELECT user_id FROM public.user WHERE user_id=(
						SELECT user_id FROM public.order_parcel WHERE tracking_code=$1 AND provider_code=$2
					)`,
		trackingCode, providerCode,
	)
	err := row.Scan(&userId)
	if err != nil {
		return "", internal.ErrDatabase{InternalError: err}
	}
	return userId, nil
}

func (r *repository) FindIdByTrackingCode(trackingCode string, providerCode entity.ProviderCode) (string, error) {
	var id string
	row := r.DB.QueryRow(`SELECT order_parcel_id FROM public.order_parcel WHERE tracking_code=$1 AND provider_code=$2`,
		trackingCode, providerCode,
	)
	err := row.Scan(&id)
	if err != nil {
		return "", internal.ErrDatabase{InternalError: err}
	}
	return id, nil
}

//FindOrderByIds can specific userId or userId <= 0 find all of order. Return map of order_parcel_id and record.
func (r *repository) FindOrderByIds(userId string, ids []string) (map[string]FindRecord, error) {
	mapOrderParcelIdRecord := make(map[string]FindRecord)

	var values []interface{}
	query := `
SELECT o.order_parcel_id, o.provider_code, o.price, o.payment_method, 
o.weight, o.width, o.length, o.height, o.tracking_code, o.status, o.cod_amount, o.created_at,
o_shp.id, o_shp.purchase_id, o_shp.status, o_shp.courier_code, o_shp.tracking_code,
o_shp_fla.id, o_shp_fla.sort_code, o_shp_fla.dst_code, o_shp_fla.sorting_line_code
FROM public.order_parcel o
LEFT JOIN public.origin ori
ON o.origin_id=ori.id
LEFT JOIN public.destination des
ON o.destination_id=des.id
LEFT JOIN public.order_parcel_shippop o_shp
ON o.order_parcel_id=o_shp.order_parcel_id
LEFT JOIN public.order_parcel_shippop_flash o_shp_fla
ON o_shp.id=o_shp_fla.order_parcel_shippop_id
WHERE o.deleted_at=0`
	if userId != "" {
		values = append(values, userId)
		query = fmt.Sprintf(`%v AND o.user_id=$%v`, query, len(values))
	}
	var idsPlaceHolder []string
	for _, id := range ids {
		values = append(values, id)
		idsPlaceHolder = append(idsPlaceHolder, fmt.Sprintf("$%v", len(values)))
	}
	query = fmt.Sprintf(`%v AND o.order_parcel_id IN (%v)`, query, strings.Join(idsPlaceHolder, ", "))

	//log.Printf("values: %v\n", values)
	//log.Printf("query: %v\n", query)

	rows, err := r.DB.Query(query, values...)
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	for rows.Next() {
		r := FindRecord{}
		o := struct {
			OrderParcelId interface{}
			ProviderCode  interface{}
			Price         interface{}
			PaymentMethod interface{}
			Weight        interface{}
			Width         interface{}
			Length        interface{}
			Height        interface{}
			TrackingCode  interface{}
			Status        interface{}
			CodAmount     interface{}
			CreatedAt     interface{}
		}{}
		oShp := struct {
			Id           interface{}
			PurchaseId   interface{}
			Status       interface{}
			CourierCode  interface{}
			TrackingCode interface{}
		}{}
		oShpFla := struct {
			Id              interface{}
			SortCode        interface{}
			DstCode         interface{}
			SortingLineCode interface{}
		}{}
		err := rows.Scan(
			&o.OrderParcelId, &o.ProviderCode, &o.Price, &o.PaymentMethod,
			&o.Weight, &o.Width, &o.Length, &o.Height, &o.TrackingCode, &o.Status, &o.CodAmount, &o.CreatedAt,
			&oShp.Id, &oShp.PurchaseId, &oShp.Status, &oShp.CourierCode, &oShp.TrackingCode,
			&oShpFla.Id, &oShpFla.SortCode, &oShpFla.DstCode, &oShpFla.SortingLineCode,
		)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		//Populate order_parcel
		r.OrderParcel = entity.OrderParcel{
			OrderParcelId: o.OrderParcelId.(string),
			ProviderCode:  entity.ProviderCode(o.ProviderCode.(int64)),
			Price:         str.ParseFloat64(o.Price.(string)),
			PaymentMethod: entity.OrderPaymentMethod(o.PaymentMethod.(int64)),
			ParcelShape: entity.ParcelShape{
				Weight: float32(o.Weight.(float64)),
				Width:  float32(o.Width.(float64)),
				Length: float32(o.Length.(float64)),
				Height: float32(o.Height.(float64)),
			},
			TrackingCode: o.TrackingCode.(string),
			Status:       o.Status.(string),
			CodAmount:    str.ParseFloat64(o.CodAmount.(string)),
			CreateDelete: entity.CreateDelete{
				CreatedAt: o.CreatedAt.(time.Time),
			},
		}
		//Populate order_parcel_shippop
		if oShp.Id != nil && r.OrderParcel.ProviderCode == entity.ProviderCodeShippop {
			r.OrderParcelShippop = &entity.OrderParcelShippop{
				Id:           oShp.Id.(int64),
				PurchaseId:   oShp.PurchaseId.(int64),
				Status:       oShp.Status.(string),
				CourierCode:  oShp.CourierCode.(string),
				TrackingCode: oShp.TrackingCode.(string),
			}
			//Populate order_parcel_shippop_flash
			if oShpFla.Id != nil && r.OrderParcelShippop.CourierCode == "FLE" {
				r.OrderParcelShippopFlash = &entity.OrderParcelShippopFlash{
					Id: oShpFla.Id.(int64),
					ShippopFlashSortingCode: entity.ShippopFlashSortingCode{
						SortCode:        oShpFla.SortCode.(string),
						DstCode:         oShpFla.DstCode.(string),
						SortingLineCode: oShpFla.SortingLineCode.(string),
					},
				}
			}
		}
		mapOrderParcelIdRecord[r.OrderParcel.OrderParcelId] = r
	}
	err = rows.Err()
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	return mapOrderParcelIdRecord, nil
}
