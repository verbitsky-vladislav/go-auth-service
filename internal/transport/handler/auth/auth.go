package auth

import (
	"auth-microservice/internal/model"
	"auth-microservice/internal/transport/handler/errors"
	"auth-microservice/internal/transport/handler/responses"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided information.
// This endpoint allows users to register with their email and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param userCreate body model.UserCreate true "User registration information"
// @Success 201 {object} responses.CreateUserResponse "Successfully registered user"
// @Failure 400 {object} errors.ErrorResponse "Bad Request"
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /api/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var userCreate model.UserCreate
	if err := c.BindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "Bad Request",
			},
		})
		return
	}

	userId, err := h.authService.Register(&userCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusInternalServerError,
				Error:     err.Error(),
				Message:   "Internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, responses.CreateUserResponse{
		Id:      userId,
		Message: "please check your email for verification",
		Status:  http.StatusOK,
	})
}

// Login godoc
// @Summary Login user
// @Description Log in a user with email and password.
// This endpoint allows users to log in using their registered credentials.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.UserLogin true "User login information"
// @Success 200 {object} responses.UserLoginResponse "Successfully logged in user"
// @Failure 400 {object} errors.ErrorResponse "Bad Request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /api/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var userLogin model.UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "Bad Request",
			},
		})
		return
	}

	userInfo, tokens, err := h.authService.Login(&userLogin)
	if err != nil {
		switch err.Error() {
		case "user not found", "incorrect password":
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
				Error: errors.Error{
					Timestamp: timestamppb.Now().String(),
					Status:    http.StatusUnauthorized,
					Error:     err.Error(),
					Message:   "Unauthorized",
				},
			})
		case "user is not verified. verification url was sent on email":
			c.JSON(http.StatusForbidden, errors.ErrorResponse{
				Error: errors.Error{
					Timestamp: timestamppb.Now().String(),
					Status:    http.StatusForbidden,
					Error:     err.Error(),
					Message:   "Forbidden",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
				Error: errors.Error{
					Timestamp: timestamppb.Now().String(),
					Status:    http.StatusInternalServerError,
					Error:     err.Error(),
					Message:   "Internal server error",
				},
			})
		}
		return
	}

	c.SetCookie("access_token", tokens.AccessToken, h.cfg.Jwt.ACCESS_LIFE_TIME, "/", h.cfg.Application.DOMAIN, true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, h.cfg.Jwt.REFRESH_LIFE_TIME, "/", h.cfg.Application.DOMAIN, true, true)

	c.JSON(http.StatusOK, responses.UserLoginResponse{
		Status:   http.StatusOK,
		Message:  "user successfully login",
		UserInfo: userInfo,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, responses.LogoutUserResponse{
		Status:  http.StatusOK,
		Message: "user successfully logout",
	})
}
