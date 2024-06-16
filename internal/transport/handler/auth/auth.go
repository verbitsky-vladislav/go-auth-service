package auth

import (
	"auth-microservice/internal/model"
	"auth-microservice/internal/transport/handler/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

func (h *Handler) Register(c *gin.Context) {

	var userCreate model.UserCreate
	if err := c.BindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     "Invalid request data",
				Message:   err.Error(),
			},
		})
		return
	}

	userId, err := h.userService.CreateUser(userCreate)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     "Bad request",
				Message:   err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, model.BrandSuccessResponse{
		Status:  "success",
		ID:      brandID,
		Message: "Brand successfully created",
	})
}
