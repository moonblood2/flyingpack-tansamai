package product

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg/dbpg"
	"strconv"
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

func (r *repository) Insert(p entity.Product) (entity.Product, error) {
	p.CreatedAt = time.Now()
	row := r.DB.QueryRow(`INSERT INTO public.product(
	user_id, name, price, created_at)
	VALUES ($1, $2, $3, $4) RETURNING id`, p.UserId, p.Name, p.Price, p.CreatedAt)
	err := row.Scan(&p.Id)
	if err != nil {
		return entity.Product{}, internal.ErrDatabase{InternalError: err}
	}
	return p, nil
}

func (r *repository) DoesIdsExistByUserId(ids []int32, userId string) ([]bool, error) {
	idsLen := len(ids)
	//existList is output.
	existList := make([]bool, idsLen)

	//idRange use bytes.Buffer to avoid concatenating string with + operator. + operator will generate unnecessary new string object.
	var values []interface{}
	var placeHolders []string
	//mem: memo the id that have found in ids, avoid duplicate id.
	mem := map[int32]bool{}
	//n is running number of place holder $n.
	n := 0
	for i := 0; i < idsLen; i++ {
		if !mem[ids[i]] {
			n += 1
			mem[ids[i]] = true
			values = append(values, strconv.FormatInt(int64(ids[i]), 10))
			placeHolders = append(placeHolders, fmt.Sprintf("$%v", n))
		}
	}
	values = append(values, userId)
	query := fmt.Sprintf("SELECT id FROM public.product WHERE id IN (%v) AND user_id=$%v", strings.Join(placeHolders, ", "), n+1)

	rows, err := r.DB.Query(query, values...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return existList, nil
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()

	//existIds exist id from database.
	existIds := map[int32]bool{}
	for rows.Next() {
		var id int32
		err := rows.Scan(&id)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		existIds[id] = true
	}
	err = rows.Err()
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	//Check input ids with id from database.
	for i := 0; i < idsLen; i++ {
		if existIds[ids[i]] {
			existList[i] = true
		} else {
			existList[i] = false
		}
	}

	return existList, nil
}

func (r *repository) FindAll() ([]entity.Product, error) {
	rows, err := r.DB.Query(`SELECT id, user_id, name, price, created_at FROM public.product WHERE deleted_at=0`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	var products []entity.Product
	for rows.Next() {
		p := entity.Product{}
		err := rows.Scan(&p.Id, &p.UserId, &p.Name, &p.Price, &p.CreatedAt)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		products = append(products, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	return products, nil
}

func (r *repository) FindByContactId(userId string) ([]entity.Product, error) {
	rows, err := r.DB.Query(`SELECT id, user_id, name, price, created_at FROM public.product WHERE user_id=$1 AND deleted_at=0`, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	var products []entity.Product
	for rows.Next() {
		p := entity.Product{}
		err := rows.Scan(&p.Id, &p.UserId, &p.Name, &p.Price, &p.CreatedAt)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		products = append(products, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	return products, nil
}

func (r *repository) FindByContactIdAndId(userId string, id int32) (entity.Product, error) {
	rows := r.DB.QueryRow(`SELECT id, user_id, name, price, created_at FROM public.product WHERE user_id=$1 AND id=$2 AND deleted_at=0`, userId, id)
	var p entity.Product
	err := rows.Scan(&p.Id, &p.UserId, &p.Name, &p.Price, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, internal.ErrNotFound{InternalError: err, Details: fmt.Sprintf("Product id %v not found. ", id)}
		}
		return entity.Product{}, internal.ErrDatabase{InternalError: err}
	}
	return p, nil
}

func (r *repository) Update(product entity.Product) (entity.Product, error) {
	var columns []string
	var values []interface{}
	//Set columns=values.
	if product.Name != "" {
		columns = append(columns, "name")
		values = append(values, product.Name)
	}
	if product.Price != 0 {
		columns = append(columns, "price")
		values = append(values, product.Price)
	}
	//Where values.
	values = append(values, product.Id)
	values = append(values, product.UserId)

	n := len(columns)
	query := dbpg.PrepareUpdateRow("public.product", columns)
	row := r.DB.QueryRow(fmt.Sprintf(`%v WHERE id=$%v AND user_id=$%v AND deleted_at=0 RETURNING id, user_id, name, price, created_at`, query, n+1, n+2), values...)
	err := row.Scan(&product.Id, &product.UserId, &product.Name, &product.Price, &product.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, internal.ErrNotFound{InternalError: err, Details: fmt.Sprintf("Product id %v not found. ", product.Id)}
		}
		return entity.Product{}, internal.ErrDatabase{InternalError: err}
	}
	return product, nil
}

func (r *repository) Delete(product entity.Product) error {
	row := r.DB.QueryRow(`UPDATE public.product SET deleted_at=$1 WHERE id=$2 AND user_id=$3 AND deleted_at=0 RETURNING id`, time.Now().Unix(), product.Id, product.UserId)
	err := row.Scan(&product.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.ErrNotFound{InternalError: err, Details: fmt.Sprintf("Product id %v not found. ", product.Id)}
		}
		return internal.ErrDatabase{InternalError: err}
	}
	return nil
}
