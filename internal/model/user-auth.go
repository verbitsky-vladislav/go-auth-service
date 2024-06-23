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

type UserGoogleInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type UserYandexInfo struct {
	Id           string   `json:"id"`
	Login        string   `json:"login"`
	ClientID     string   `json:"client_id"`
	DisplayName  string   `json:"display_name"`
	RealName     string   `json:"real_name"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Sex          string   `json:"sex"`
	DefaultEmail string   `json:"default_email"`
	Emails       []string `json:"emails"`
	AvatarID     string   `json:"default_avatar_id"`
	Birthday     string   `json:"birthday"`
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
