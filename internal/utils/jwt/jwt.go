package jwt

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/model"
	"auth-microservice/internal/service"
	"auth-microservice/pkg/logger"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type jwtService struct {
	cfg *config.Config
}

func NewJwtService(cfg *config.Config) service.JwtService {
	return &jwtService{
		cfg: cfg,
	}
}

func (j *jwtService) GenerateTokens(user model.UserInfo) (model.Tokens, error) {
	var tokens model.Tokens

	accessToken, err := j.generateAccessToken(user)
	if err != nil {
		return tokens, err
	}
	tokens.AccessToken = accessToken

	refreshToken, err := j.GenerateRefreshToken()
	if err != nil {
		return tokens, err
	}
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func (j *jwtService) generateAccessToken(user model.UserInfo) (string, error) {
	claims := model.UserClaims{
		User: user,
		Claims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(j.cfg.Jwt.ACCESS_LIFE_TIME) * time.Second).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.cfg.Jwt.SECRET_KEY))
	if err != nil {
		return "", logger.Error(err, "error signing access token")
	}

	return signedToken, nil
}

func (j *jwtService) VerifyAccessToken(tokenString string) (model.UserInfo, error) {
	claims := &model.UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.Jwt.SECRET_KEY), nil
	})
	if err != nil {
		return model.UserInfo{}, logger.Error(err, "error parsing access token")
	}
	if !token.Valid {
		return model.UserInfo{}, fmt.Errorf("invalid access token")
	}

	return claims.User, nil
}

func (j *jwtService) GenerateRefreshToken() (string, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Duration(j.cfg.Jwt.REFRESH_LIFE_TIME) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.cfg.Jwt.SECRET_KEY))
	if err != nil {
		return "", logger.Error(err, "error signing refresh token")
	}

	return signedToken, nil
}

func (j *jwtService) VerifyRefreshToken(tokenString string) error {
	claims := jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.Jwt.SECRET_KEY), nil
	})
	if err != nil {
		return logger.Error(err, "error parsing refresh token")
	}
	if !token.Valid {
		return fmt.Errorf("invalid refresh token")
	}

	return nil
}
