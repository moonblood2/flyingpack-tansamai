package user

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/entity"
)

// repository implement Repository interface.
type repository struct {
	DB *sql.DB
}

// NewPostgresRepository return new repository that implement by PostgreSQL.
func NewPostgresRepository(db *sql.DB) *repository {
	return &repository{
		DB: db,
	}
}

// DoesEmailExist return true if exist.
func (r *repository) DoesEmailExist(email string) (bool, error) {
	var userId string
	err := r.DB.QueryRow(`SELECT user_id FROM public."user" WHERE email = $1 AND deleted_at=0`, email).Scan(&userId)
	if err != nil {
		// If an err is no rows, e-mail doesn't exists, return false.
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, internal.ErrDatabase{InternalError: err}
	}
	return true, nil
}

// Create create user and return user id.
func (r *repository) Create(user entity.User) (entity.User, error) {
	//Insert user
	row := r.DB.QueryRow(
		`INSERT INTO public."user" (user_id, name, email, role, password, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`,
		user.Id, user.Name, user.Email, user.Role, user.Password, time.Now(),
	)
	err := row.Scan(&user.Id)
	if err != nil {
		return entity.User{}, internal.ErrDatabase{InternalError: err}
	}
	return user, nil
}

// FindAll return all users.
func (r *repository) FindAll() ([]entity.User, error) {
	//TODO 1, Pagination users.
	var users []entity.User
	rows, err := r.DB.Query(`SELECT user_id, email, name, role FROM public."user" WHERE deleted_at=0`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, internal.ErrNotFound{Details: "Users not found."}
		}
		return users, internal.ErrDatabase{InternalError: err}
	}
	defer rows.Close()
	for rows.Next() {
		user := entity.User{}
		if err := rows.Scan(&user.Id, &user.Email, &user.Name, &user.Role); err != nil {
			return users, internal.ErrDatabase{InternalError: err}
		}
		users = append(users, user)
	}
	return users, nil
}

// FindById return user for specific id.
func (r *repository) FindById(id string) (entity.User, error) {
	var u entity.User
	row := r.DB.QueryRow(`SELECT user_id, email, name, role FROM public."user" WHERE user_id = $1 AND deleted_at=0`, id)
	err := row.Scan(&u.Id, &u.Email, &u.Name, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, internal.ErrNotFound{Details: fmt.Sprintf("User %v not found.", id)}
		}
		return u, internal.ErrDatabase{InternalError: err}
	}
	return u, nil
}

// FindByEmail return user for specific e-mail.
func (r *repository) FindByEmail(email string) (entity.User, error) {
	u := entity.User{}
	row := r.DB.QueryRow(`SELECT user_id, email, name, role, password FROM public."user" WHERE email = $1 AND deleted_at=0`, email)
	err := row.Scan(&u.Id, &u.Email, &u.Name, &u.Role, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, internal.ErrNotFound{Details: fmt.Sprintf("Email %v not found.", email)}
		}
		return entity.User{}, internal.ErrDatabase{InternalError: err}
	}
	return u, nil
}

// Update update user by id in specific filed.
func (r *repository) Update(user entity.User) (entity.User, error) {
	_, err := r.DB.Exec(`UPDATE "public".user SET name=$1, role=$2, password=$3 WHERE user_id=$4`, user.Name, user.Role, user.Password, user.Id)
	if err != nil {
		return entity.User{}, internal.ErrDatabase{InternalError: err}
	}
	return user, nil
}

// Delete delete user by id, use soft-delete set deleted_at to current time.
func (r *repository) Delete(id string) error {
	_, err := r.DB.Exec(`UPDATE "public".user SET deleted_at=$1 WHERE user_id=$2`, time.Now(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.ErrNotFound{Details: fmt.Sprintf("User %v not found.", id)}
		}
		return internal.ErrDatabase{InternalError: err}
	}
	return nil
}
