package main

import (
	_ "auth-microservice/docs"
	"auth-microservice/internal/config"
	"auth-microservice/internal/repository/postgres"
	"auth-microservice/internal/service/auth"
	"auth-microservice/internal/service/auth/social/google"
	"auth-microservice/internal/service/auth/social/vk"
	"auth-microservice/internal/service/auth/social/yandex"
	"auth-microservice/internal/service/cache"
	"auth-microservice/internal/service/user"
	"auth-microservice/internal/transport/handler"
	"auth-microservice/internal/utils/jwt"
	"auth-microservice/internal/utils/verification"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/mailer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	cfg := config.New()

	db, err := postgres.NewDB(
		cfg.Database.HOST,
		cfg.Database.PORT,
		cfg.Database.USER,
		cfg.Database.PASSWORD,
		cfg.Database.DATABASE,
	)
	if err != nil {
		logger.Panic(err, "failed to connect to database")
	}

	redisClient, err := cache.NewRedisCache(
		cfg.Redis.HOST+":"+cfg.Redis.PORT,
		cfg.Redis.PASSWORD,
		0,
	)
	if err != nil {
		logger.Panic(err, "failed to create redis client")
	}

	userRepo := postgres.NewUserRepository(db)

	mailerService := mailer.NewMailer(cfg.Mailer.USERNAME, cfg.Mailer.PASSWORD)
	userService := user.NewUserService(userRepo)
	verificationService := verification.NewVerificationService(cfg)
	jwtService := jwt.NewJwtService(cfg)
	authService := auth.NewAuthService(
		userService,
		redisClient,
		mailerService,
		verificationService,
		jwtService,
	)
	googleService := google.NewGoogleService(cfg, userService)
	vkService := vk.NewVkService(cfg, userService)
	yandexService := yandex.NewYandexService(cfg, userService)

	_, s := handler.NewHandler(
		cfg,
		userService,
		authService,
		googleService,
		vkService,
		yandexService,
		jwtService,
	)
	router := gin.Default()

	//router.Use(middleware.CORSMiddleware())
	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5050", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.RegisterRoutes(router)

	if err := router.Run(cfg.Application.PORT); err != nil {
		logger.Panic(err, "failed to start server")
	}
}
