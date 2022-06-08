package origin

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
	row := r.DB.QueryRow(`SELECT phone_number FROM public.origin WHERE phone_number=$1 AND deleted_at=0`, phoneNumber)
	err := row.Scan(&phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, internal.ErrDatabase{InternalError: err}
	}
	return true, nil
}

func (r *repository) DoesPhoneNumbersExist(phoneNumbers []string) ([]bool, error) {
	pLen := len(phoneNumbers)
	//Initialize existList, all elements are false.
	existList := make([]bool, pLen)
	//Make unique phone number list by mem.
	mem := make(map[string]bool)
	for _, v := range phoneNumbers {
		if _, ok := mem[v]; !ok {
			mem[v] = true
		}
	}
	//Prepare placeholder list and value list for query.
	memLen := len(mem)
	valuesPlaceholder := make([]string, memLen)
	values := make([]interface{}, memLen)
	n := 0
	for k := range mem {
		valuesPlaceholder[n] = fmt.Sprintf("$%v", n+1)
		values[n] = k
		n++
	}
	//Make query.
	query := fmt.Sprintf(`SELECT phone_number FROM public."origin" WHERE phone_number IN (%v)`, strings.Join(valuesPlaceholder, ", "))
	fmt.Printf("query: %v\n", query)
	rows, err := r.DB.Query(query, values...)
	if err != nil {
		//If not found any rows return all false.
		if errors.Is(err, sql.ErrNoRows) {
			return existList, nil
		}
		return nil, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	//Mem found phone number.
	found := make(map[string]bool)
	for rows.Next() {
		var pn string
		err := rows.Scan(&pn)
		if err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
		found[pn] = true
	}
	//For each phone number and check in found mem.
	for i, pn := range phoneNumbers {
		//If found set to true.
		if _, ok := found[pn]; ok {
			existList[i] = true
		}
	}
	return existList, nil
}

func (r *repository) UpsertByPhoneNumber(origins []entity.Origin) ([]entity.Origin, error) {
	//Group each origin by phone number. key = phone_number, value = list of origins index.
	group := make(map[string][]int)
	for i, v := range origins {
		group[v.PhoneNumber] = append(group[v.PhoneNumber], i)
	}
	//Unique origin.
	unqOrgLen := len(group)
	//Use last origin in each group.
	unqOrg := make([]entity.Origin, unqOrgLen)
	i := 0
	for _, v := range group {
		unqOrg[i] = origins[v[len(v)-1]]
		i++
	}
	//Make value placeholder list and value list.
	vpl := make([]string, unqOrgLen)
	var values []interface{}
	n := 0
	for i := 0; i < unqOrgLen; i++ {
		vpl[i] = fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v)", n+1, n+2, n+3, n+4, n+5, n+6, n+7)
		n += 7

		values = append(values, unqOrg[i].Name)
		values = append(values, unqOrg[i].PhoneNumber)
		values = append(values, unqOrg[i].Address)
		values = append(values, unqOrg[i].District)
		values = append(values, unqOrg[i].State)
		values = append(values, unqOrg[i].Province)
		values = append(values, unqOrg[i].Postcode)
	}
	//Make query
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, internal.ErrDatabase{InternalError: err}
	}
	//Set current id to max id in column.
	if _, err := tx.Exec(`SELECT setval('public.origin_id_seq', MAX(id)) FROM public.origin;`); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, internal.ErrDatabase{InternalError: err}
		}
	}
	//Query upsert.
	query := fmt.Sprintf(`INSERT INTO public.origin (name, phone_number, address, district, state, province, postcode) VALUES %v ON CONFLICT ON CONSTRAINT origin_phone_number_deleted_at_key DO UPDATE SET name=EXCLUDED.name, address=EXCLUDED.address, district=EXCLUDED.district, state=EXCLUDED.state, province=EXCLUDED.province, postcode=EXCLUDED.postcode RETURNING id;`, strings.Join(vpl, ", "))
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
		if err := rows.Scan(&unqOrg[i].Id); err != nil {
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
	//O(nlog(n)), map unique origin to output.
	for _, u := range unqOrg {
		for _, i := range group[u.PhoneNumber] {
			origins[i] = u
		}
	}
	return origins, nil
}
