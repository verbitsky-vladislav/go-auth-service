package model

import "time"

type User struct {
	Id               *string    `json:"id" db:"id"`
	Username         *string    `json:"username" db:"username"`
	Email            *string    `json:"email" db:"email"`
	Phone            *string    `json:"phone" db:"phone"`
	Password         *string    `json:"password" db:"password"`
	IsVerified       bool       `json:"is_verified" db:"is_verified"`
	GoogleAuthSecret *string    `json:"google_auth_secret" db:"google_auth_secret"`
	CreatedAt        *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at" db:"updated_at"`
}

type UserCreate struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Password string `json:"password" db:"password"`
}

type UserCreateFromGoogle struct {
	Id         string `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	Email      string `json:"email" db:"email"`
	IsVerified bool   `json:"is_verified" db:"is_verified"`
	Logo       string `json:"logo" db:"logo"`
}

type UserCreateFromYandex struct {
	Id           string `json:"id"`
	RealName     string `json:"real_name"`
	IsVerified   bool   `json:"is_verified"`
	DefaultEmail string `json:"default_email"`
	Avatar       string `json:"default_avatar"`
}

type UserUpdate struct {
	Username         string `json:"username" db:"username"`
	Email            string `json:"email" db:"email"`
	Phone            string `json:"phone" db:"phone"`
	GoogleAuthSecret string `json:"google_auth_secret" db:"google_auth_secret"`
	IsVerified       bool   `json:"is_verified" db:"is_verified"`
}
