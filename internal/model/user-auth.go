package model

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserClaims struct {
	User   UserInfo
	Claims jwt.StandardClaims
}

func (u *UserClaims) Valid() error {
	// check if token is expired
	if u.Claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}
	// TODO: add additional validation

	return nil
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
