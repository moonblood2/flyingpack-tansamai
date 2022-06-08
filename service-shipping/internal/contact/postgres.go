package contact

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

func (r *repository) DoesPhoneNumberExist(phoneNumber string) (bool, error) {
	row := r.DB.QueryRow(`SELECT phone_number FROM public.contact WHERE phone_number=$1 AND deleted_at=0`, phoneNumber)
	err := row.Scan(&phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, internal.ErrDatabase{InternalError: err}
	}
	return true, nil
}

func (r *repository) FindByUserId(userId string) (entity.Contact, error) {
	s := entity.Contact{}
	row := r.DB.QueryRow(`SELECT id, user_id, name, phone_number, address, district, state, province, postcode FROM public.contact WHERE user_id=$1 AND deleted_at=0`, userId)
	err := row.Scan(&s.Id, &s.UserId, &s.Name, &s.PhoneNumber, &s.Address, &s.District, &s.State, &s.Province, &s.Postcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Contact{}, internal.ErrNotFound{InternalError: err, Details: fmt.Sprintf("contact for user id %v not found", userId)}
		}
		return entity.Contact{}, internal.ErrDatabase{InternalError: err}
	}
	return s, nil
}

func (r *repository) Create(s entity.Contact) (entity.Contact, error) {
	row := r.DB.QueryRow(`INSERT INTO public.contact(
	user_id, name, phone_number, address, district, state, province, postcode, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		s.UserId, s.Name, s.PhoneNumber, s.Address, s.District, s.State, s.Province, s.Postcode, time.Now(),
	)
	if err := row.Scan(&s.Id); err != nil {
		return entity.Contact{}, internal.ErrDatabase{InternalError: err}
	}
	return s, nil
}

func (r *repository) CreateWithUser(user entity.User, contact entity.Contact) (entity.User, entity.Contact, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return entity.User{}, entity.Contact{}, internal.ErrDatabase{InternalError: err}
	}
	//Insert user.
	user.CreatedAt = time.Now()
	row := tx.QueryRow(`INSERT INTO public."user"(
	user_id, email, name, role, password, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		user.Id, user.Email, user.Name, user.Role, user.Password, user.CreatedAt,
	)
	//Insert contact.
	contact.UserId = user.Id
	contact.CreatedAt = time.Now()
	row = tx.QueryRow(`INSERT INTO public.contact(
	user_id, name, phone_number, address, district, state, province, postcode, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		contact.UserId, contact.Name, contact.PhoneNumber, contact.Address, contact.District, contact.State, contact.Province, contact.Postcode, contact.CreatedAt,
	)
	err = row.Scan(&contact.Id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return entity.User{}, entity.Contact{}, internal.ErrDatabase{InternalError: err}
		}
	}
	err = tx.Commit()
	if err != nil {
		return entity.User{}, entity.Contact{}, internal.ErrDatabase{InternalError: err}
	}
	return user, contact, nil
}
