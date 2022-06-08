package destination

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"strings"
)

type repository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) DoesPhoneNumberExist(phoneNumber string) (bool, error) {
	row := r.DB.QueryRow(`SELECT phone_number FROM public.destination WHERE phone_number=$1 AND deleted_at=0`, phoneNumber)
	err := row.Scan(&phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, internal.ErrDatabase{InternalError: err}
	}
	return true, nil
}

func (r *repository) UpsertByPhoneNumber(destinations []entity.Destination) ([]entity.Destination, error) {
	//Group each destination by phone number. key = phone_number, value = list of destinations index.
	group := make(map[string][]int)
	for i, v := range destinations {
		group[v.PhoneNumber] = append(group[v.PhoneNumber], i)
	}
	unqDesLen := len(group)
	//Use last destination in each group.
	unqDes := make([]entity.Destination, unqDesLen)
	i := 0
	for _, v := range group {
		unqDes[i] = destinations[v[len(v)-1]]
		i++
	}
	//Make value placeholder list and value list.
	vpl := make([]string, unqDesLen)
	var values []interface{}
	n := 0
	for i := 0; i < unqDesLen; i++ {
		vpl[i] = fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v)", n+1, n+2, n+3, n+4, n+5, n+6, n+7)
		n += 7

		values = append(values, unqDes[i].Name)
		values = append(values, unqDes[i].PhoneNumber)
		values = append(values, unqDes[i].Address)
		values = append(values, unqDes[i].District)
		values = append(values, unqDes[i].State)
		values = append(values, unqDes[i].Province)
		values = append(values, unqDes[i].Postcode)
	}
	//Make query
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	//Set current id to max id in column.
	if _, err := tx.Exec(`SELECT setval('public.destination_id_seq', MAX(id)) FROM public.destination;`); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
	}
	//Query upsert.
	query := fmt.Sprintf(`INSERT INTO public.destination (name, phone_number, address, district, state, province, postcode) VALUES %v ON CONFLICT ON CONSTRAINT destination_phone_number_deleted_at_key DO UPDATE SET name=EXCLUDED.name, address=EXCLUDED.address, district=EXCLUDED.district, state=EXCLUDED.state, province=EXCLUDED.province, postcode=EXCLUDED.postcode RETURNING id;`, strings.Join(vpl, ", "))
	rows, err := tx.Query(query, values...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	//Retrieve id.
	i = 0
	for rows.Next() {
		if err := rows.Scan(&unqDes[i].Id); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, internal.ErrDatabase{InternalError: err}
			}
			return nil, internal.ErrDatabase{InternalError: err}
		}
		i++
	}
	if err := rows.Err(); err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	//Commit.
	if err := tx.Commit(); err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	//O(nlog(n)), map unique destination to output.
	for _, u := range unqDes {
		for _, i := range group[u.PhoneNumber] {
			destinations[i] = u
		}
	}
	return destinations, nil
}
