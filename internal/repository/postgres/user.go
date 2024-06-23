package postgres

import (
	"auth-microservice/internal/model"
	"auth-microservice/pkg/logger"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *model.UserCreate) (string, error) {
	query := `INSERT INTO users (email, phone, password, username, created_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := r.db.QueryRow(query, user.Email, user.Phone, user.Password, user.Username, time.Now()).Scan(&id)
	if err != nil {
		return "", logger.Error(err, "failed to create user")
	}
	return id, nil
}

func (r *UserRepository) Update(id string, user *model.UserUpdate) error {
	query := "UPDATE users SET"
	params := make([]interface{}, 0)
	idx := 1
	changes := false

	if user.Phone != "" {
		query += " phone = $" + strconv.Itoa(idx) + ","
		params = append(params, user.Phone)
		idx++
		changes = true
	}

	if user.Email != "" {
		query += " email = $" + strconv.Itoa(idx) + ","
		params = append(params, user.Email)
		idx++
		changes = true
	}

	if user.GoogleAuthSecret != "" {
		query += " google_auth_secret = $" + strconv.Itoa(idx) + ","
		params = append(params, user.GoogleAuthSecret)
	}

	if user.IsVerified != false {
		query += " is_verified = $" + strconv.Itoa(idx) + ","
		params = append(params, user.IsVerified)
		idx++
		changes = true
	}

	if !changes {
		return logger.Error(errors.New(query), "nothing to update user")
	}

	query = strings.TrimSuffix(query, ",")
	query += " WHERE id = $" + strconv.Itoa(idx)
	params = append(params, id)

	_, err := r.db.Exec(query, params...)
	if err != nil {
		return logger.Error(err, "failed to update user")
	}
	return nil
}

func (r *UserRepository) FindById(id string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, logger.Error(err, "failed to find user")
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, logger.Error(err, "failed to find user")
	}
	return &user, nil
}
