package auth

//import (
//	"auth-microservice/internal/model"
//	"auth-microservice/internal/transport/handler/errors"
//	"github.com/gin-gonic/gin"
//	"google.golang.org/protobuf/types/known/timestamppb"
//	"net/http"
//)
//
//func (h *Handler) My(c *gin.Context) {
//
//	//var brand model.UserCreate
//	//if err := c.BindJSON(&brand); err != nil {
//	//	c.JSON(http.StatusBadRequest, errors.ErrorResponse{
//	//		Error: errors.Error{
//	//			Timestamp: timestamppb.Now().String(),
//	//			Status:    http.StatusBadRequest,
//	//			Error:     "Invalid request data",
//	//			Message:   err.Error(),
//	//		},
//	//	})
//	//	return
//	//}
//	//
//	//brandID, err := h.brandService.CreateBrand(&brand)
//	//if err != nil {
//	//	c.JSON(http.StatusBadRequest, model.ErrorResponse{
//	//		Error: model.Error{
//	//			Timestamp: timestamppb.Now().String(),
//	//			Status:    http.StatusBadRequest,
//	//			Error:     "Bad request",
//	//			Message:   err.Error(),
//	//		},
//	//	})
//	//	return
//	}
//
//	c.JSON(http.StatusCreated, model.BrandSuccessResponse{
//		Status:  "success",
//		ID:      brandID,
//		Message: "Brand successfully created",
//	})
//}
