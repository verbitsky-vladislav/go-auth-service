package auth

import (
	"auth-microservice/internal/transport/handler/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

// Confirm godoc
// @Summary      Confirm email verification
// @Description  Confirms a user's email using a verification token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        verification_token  path      string  true  "Verification Token"
// @Success      302  {string}       string    "Redirect to frontend URL on success"
// @Failure      400  {object}       errors.ErrorResponse  "Verification token invalid"
// @Router       /confirm/verify/{verification_token} [get]
func (h *Handler) Confirm(c *gin.Context) {
	verificationToken := c.Param("verification_token")

	_, err := h.authService.ConfirmEmail(verificationToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "Verification token invalid",
			},
		})
	}

	c.Redirect(http.StatusFound, h.cfg.Application.FRONTEND_URL)

}
