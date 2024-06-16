package handler

import (
	"github.com/gin-gonic/gin"
)

type Services struct {
}

type Handlers struct {
}

// NewHandler создает экземпляр Handler с предоставленными сервисами
func NewHandler() (*Services, *Handlers) {
	return &Services{}, &Handlers{}
}

func (h *Handlers) RegisterRoutes(router *gin.Engine) {
	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("/register")
		authRouter.POST("/login")
		authRouter.POST("/logout")

		router.GET("/my")

		// для верефикации пользователя после регистрации
		authRouter.POST("/confirm/send")
		authRouter.POST("/confirm/verify/:verification_token")
	}

	//otpRouter := router.Group("/api/auth/otp")
	//{
	//	otpRouter.POST("/send")
	//	otpRouter.POST("/verify")
	//}
	//
	//socialRouter := router.Group("/api/auth/social")
	//{
	//	socialRouter.GET("/facebook")
	//	socialRouter.GET("/instagram")
	//	socialRouter.GET("/git")
	//	socialRouter.GET("/google")
	//}
	//
	//secureRouter := router.Group("/api/secure")
	//{
	//	secureRouter.POST("/reset-password/request")
	//	secureRouter.POST("/reset-password/verify")
	//	secureRouter.POST("/change-password")
	//
	//	secureRouter.POST("/change-email/request")
	//	secureRouter.POST("/change-email/verify")
	//}
}
