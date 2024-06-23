package handler

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/service"
	"auth-microservice/internal/transport/handler/auth"
	"auth-microservice/internal/transport/handler/common/middleware"
	"auth-microservice/internal/transport/handler/social/google"
	"auth-microservice/internal/utils/jwt"
	"github.com/gin-gonic/gin"
)

type Services struct {
	cfg           *config.Config
	userService   service.UserService
	authService   service.AuthService
	googleService service.GoogleService
	jwtService    service.JwtService
}

type Handlers struct {
	cfg           *config.Config
	authHandler   auth.Handler
	googleHandler google.Handler
}

func NewHandler(
	cfg *config.Config,
	userService service.UserService,
	authService service.AuthService,
	googleService service.GoogleService,
	jwtService service.JwtService,
) (*Services, *Handlers) {
	return &Services{
			cfg:           cfg,
			userService:   userService,
			authService:   authService,
			googleService: googleService,
		}, &Handlers{
			cfg:           cfg,
			authHandler:   *auth.NewAuthHandler(cfg, userService, authService),
			googleHandler: *google.NewGoogleHandler(cfg, googleService, userService, jwtService),
		}
}

func (h *Handlers) RegisterRoutes(router *gin.Engine) {

	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("/register", h.authHandler.Register)
		authRouter.POST("/login", h.authHandler.Login)
		authRouter.POST("/logout", h.authHandler.Logout)

		authRouter.GET("/my", middleware.AuthMiddleware(jwt.NewJwtService(h.cfg), h.cfg), h.authHandler.My) // todo сделать

		authRouter.GET("/confirm/verify/:verification_token", h.authHandler.Confirm)
	}

	//otpRouter := router.Group("/api/auth/otp")
	//{
	//	otpRouter.POST("/send")
	//	otpRouter.POST("/verify")
	//}
	//
	socialRouter := router.Group("/api/google")
	{
		//googleRouter := socialRouter.Group("google")
		//{
		socialRouter.GET("login", h.googleHandler.GoogleLogin)
		socialRouter.GET("callback", h.googleHandler.GoogleCallback)
		//}
		//socialRouter.GET("/instagram")
		//socialRouter.GET("/git")
		//socialRouter.GET("/google")
	}
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
