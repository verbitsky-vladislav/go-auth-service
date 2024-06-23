package auth

import (
	"auth-microservice/internal/transport/handler/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

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

	//c.JSON(http.StatusOK, responses.VerificatedUserResponse{
	//	Status:  http.StatusOK,
	//	Message: "user verified successfully",
	//	//UserInfo: userInfo,
	//})
	c.Redirect(http.StatusFound, h.cfg.Application.FRONTEND_URL)

}
