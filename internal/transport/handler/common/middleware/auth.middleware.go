package middleware

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/service"
	"auth-microservice/internal/transport/handler/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

func AuthMiddleware(jwtService service.JwtService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
				Error: errors.Error{
					Timestamp: timestamppb.Now().String(),
					Status:    http.StatusUnauthorized,
					Error:     err.Error(),
					Message:   "access_token cookie not found",
				},
			})
			c.Abort()
			return
		}

		userInfo, err := jwtService.VerifyAccessToken(accessToken)
		if err != nil {
			// Проверяем, истек ли срок действия access token
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					// Access token истек, попытаемся обновить токены
					refreshToken, err := c.Cookie("refresh_token")
					if err != nil {
						c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
							Error: errors.Error{
								Timestamp: timestamppb.Now().String(),
								Status:    http.StatusUnauthorized,
								Error:     err.Error(),
								Message:   "refresh_token cookie not found",
							},
						})
						c.Abort()
						return
					}

					newTokens, err := jwtService.RefreshTokens(refreshToken)
					if err != nil {
						c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
							Error: errors.Error{
								Timestamp: timestamppb.Now().String(),
								Status:    http.StatusUnauthorized,
								Error:     err.Error(),
								Message:   "refresh token verify failed",
							},
						})
						c.Abort()
						return
					}

					c.SetCookie("access_token", newTokens.AccessToken, cfg.Jwt.ACCESS_LIFE_TIME, "/", cfg.Application.DOMAIN, true, true)
					c.SetCookie("refresh_token", newTokens.RefreshToken, cfg.Jwt.REFRESH_LIFE_TIME, "/", cfg.Application.DOMAIN, true, true)

					userInfo, err = jwtService.VerifyAccessToken(newTokens.AccessToken)
					if err != nil {
						c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
							Error: errors.Error{
								Timestamp: timestamppb.Now().String(),
								Status:    http.StatusUnauthorized,
								Error:     err.Error(),
								Message:   "access_token verify failed after refresh",
							},
						})
						c.Abort()
						return
					}
				}
			} else {
				c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
					Error: errors.Error{
						Timestamp: timestamppb.Now().String(),
						Status:    http.StatusUnauthorized,
						Error:     err.Error(),
						Message:   "access_token verify failed",
					},
				})
				c.Abort()
				return
			}
		}
		c.Set("user", userInfo)
		c.Next()
	}
}
