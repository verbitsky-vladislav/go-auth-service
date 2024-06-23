package auth

import (
	"auth-microservice/internal/model"
	"auth-microservice/internal/transport/handler/errors"
	"auth-microservice/internal/transport/handler/responses/base.auth.responses"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

// My godoc
// @Summary Get current user information
// @Description Retrieves information about the currently authenticated user.
// This endpoint requires a valid authentication token obtained from cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} base_auth_responses.UserMyResponse "User information"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 403 {object} errors.ErrorResponse "Forbidden"
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /api/auth/my [get]
func (h *Handler) My(c *gin.Context) {
	userInfoContext, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusInternalServerError,
				Error:     "error while getting user info",
				Message:   "Status internal error",
			},
		})
		return
	}

	userInfo, ok := userInfoContext.(model.UserInfo)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusInternalServerError,
				Error:     "error while converting user info to type",
				Message:   "Status internal error",
			},
		})
		return
	}

	user, err := h.userService.FindUserById(userInfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusInternalServerError,
				Error:     err.Error(),
				Message:   "Status internal error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, base_auth_responses.UserMyResponse{
		UserInfo: user,
		Message:  "user info successfully retrieved",
		Status:   http.StatusOK,
	})
	return
}
