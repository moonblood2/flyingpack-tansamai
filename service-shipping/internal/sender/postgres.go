package sender

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
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

func (r *repository) DoesIdExist(id int64) (bool, error) {
	row := r.DB.QueryRow(`SELECT id FROM public.sender WHERE id=$1 AND deleted_at=0`, id)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, internal.ErrDatabase{InternalError: err}
	}
	return true, nil
}

func (r *repository) DoesPhoneNumberExist(phoneNumber string) (bool, error) {
	row := r.DB.QueryRow(`SELECT phone_number FROM public.sender WHERE phone_number=$1 AND deleted_at=0`, phoneNumber)
	err := row.Scan(&phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, internal.ErrDatabase{InternalError: err}
	}
	return true, nil
}

func (r *repository) FindById(id int64) (entity.Sender, error) {
	s := entity.Sender{}
	row := r.DB.QueryRow(`SELECT id, sender_type, national_id_number, tax_id_number, passport_number, birth_date, name, phone_number, address, district, state, province, postcode
	FROM public.sender WHERE id=$1 AND deleted_at=0`, id)
	err := row.Scan(&s.Id, &s.SenderType, &s.NationalIdNumber, &s.TaxIdNumber, &s.PassportNumber, &s.BirthDate, &s.Name, &s.PhoneNumber, &s.Address, &s.District, &s.State, &s.Province, &s.Postcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Sender{}, internal.ErrNotFound{InternalError: err, Details: fmt.Sprintf("sender id %v not found.", id)}
		}
		return entity.Sender{}, internal.ErrDatabase{InternalError: err}
	}
	return s, nil
}

func (r *repository) FindByPhoneNumber(phoneNumber string) (entity.Sender, error) {
	s := entity.Sender{}
	row := r.DB.QueryRow(`SELECT id, sender_type, national_id_number, tax_id_number, passport_number, birth_date, name, phone_number, address, district, state, province, postcode
	FROM public.sender WHERE phone_number=$1 AND deleted_at=0`, phoneNumber)
	err := row.Scan(&s.Id, &s.SenderType, &s.NationalIdNumber, &s.TaxIdNumber, &s.PassportNumber, &s.BirthDate, &s.Name, &s.PhoneNumber, &s.Address, &s.District, &s.State, &s.Province, &s.Postcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Sender{}, internal.ErrNotFound{InternalError: err, Details: fmt.Sprintf("phone_number %v not found.", phoneNumber)}
		}
		return entity.Sender{}, internal.ErrDatabase{InternalError: err}
	}
	return s, nil
}

func (r *repository) Create(s entity.Sender) (entity.Sender, error) {
	exist, err := r.DoesPhoneNumberExist(s.PhoneNumber)
	if err != nil {
		return entity.Sender{}, err
	}
	if exist {
		return entity.Sender{}, internal.ErrDatabase{InternalError: fmt.Errorf("Sender.Create(), Phone number %v duplicate. ", s.PhoneNumber)}
	}
	row := r.DB.QueryRow(`INSERT INTO public.sender(
	sender_type, national_id_number, tax_id_number, passport_number, birth_date, name, phone_number, address, district, state, province, postcode, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`,
		s.SenderType, s.NationalIdNumber, s.TaxIdNumber, s.PassportNumber, s.BirthDate, s.Name, s.PhoneNumber, s.Address, s.District, s.State, s.Province, s.Postcode, time.Now(),
	)
	if err := row.Scan(&s.Id); err != nil {
		return entity.Sender{}, internal.ErrDatabase{InternalError: err}
	}
	return s, nil
}

func (r *repository) UpdateById(s entity.Sender, id int64) (entity.Sender, error) {
	//Find previous sender.
	ps, err := r.FindById(id)
	if err != nil {
		return entity.Sender{}, err
	}
	//Populate id.
	s.Id = ps.Id
	//Check if input sender is empty field, update with previous sender.
	if s.SenderType == 0 {
		s.SenderType = ps.SenderType
	}
	if s.NationalIdNumber == "" {
		s.NationalIdNumber = ps.NationalIdNumber
	}
	if s.TaxIdNumber == "" {
		s.TaxIdNumber = ps.TaxIdNumber
	}
	if s.PassportNumber == "" {
		s.PassportNumber = ps.PassportNumber
	}
	if s.BirthDate == "" {
		s.BirthDate = ps.BirthDate
	}
	if s.Name == "" {
		s.Name = ps.Name
	}
	if s.PhoneNumber == "" {
		s.PhoneNumber = ps.PhoneNumber
	}
	if s.Address == "" {
		s.Address = ps.Address
	}
	if s.District == "" {
		s.District = ps.District
	}
	if s.State == "" {
		s.State = ps.State
	}
	if s.Province == "" {
		s.Province = ps.Province
	}
	if s.Postcode == "" {
		s.Postcode = ps.Postcode
	}
	_, err = r.DB.Exec(`UPDATE public.sender
	SET sender_type=$1, national_id_number=$2, tax_id_number=$3, passport_number=$4, birth_date=$5, name=$6, phone_number=$7, address=$8, district=$9, state=$10, province=$11, postcode=$12
	WHERE id=$13`,
		s.SenderType, s.NationalIdNumber, s.TaxIdNumber, s.PassportNumber, s.BirthDate, s.Name, s.PhoneNumber, s.Address, s.District, s.State, s.Province, s.Postcode,
		id,
	)
	if err != nil {
		return entity.Sender{}, internal.ErrDatabase{InternalError: err}
	}
	return s, nil
}

//UpdateByPhoneNumber every field in sender must have value.
func (r *repository) UpdateByPhoneNumber(s entity.Sender, phoneNumber string) (entity.Sender, error) {
	if s.Validate() != nil {
		return entity.Sender{}, nil
	}
	row := r.DB.QueryRow(`UPDATE public.sender
	SET sender_type=$1, national_id_number=$2, tax_id_number=$3, passport_number=$4, birth_date=$5, name=$6, address=$7, district=$8, state=$9, province=$10, postcode=$11
	WHERE phone_number=$12 RETURNING id`,
		s.SenderType, s.NationalIdNumber, s.TaxIdNumber, s.PassportNumber, s.BirthDate, s.Name, s.Address, s.District, s.State, s.Province, s.Postcode,
		phoneNumber,
	)
	err := row.Scan(&s.Id)
	if err != nil {
		return entity.Sender{}, internal.ErrDatabase{InternalError: err}
	}
	return s, nil
}

func (r *repository) SaveByPhoneNumber(sender entity.Sender, phoneNumber string) (entity.Sender, error) {
	exist, err := r.DoesPhoneNumberExist(phoneNumber)
	if err != nil {
		return entity.Sender{}, err
	}
	//If not exist create, if not update.
	if !exist {
		return r.Create(sender)
	}
	return r.UpdateByPhoneNumber(sender, phoneNumber)
}
