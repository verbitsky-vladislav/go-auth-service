package main

import (
	_ "auth-microservice/docs"
	"auth-microservice/internal/config"
	"auth-microservice/internal/repository/postgres"
	"auth-microservice/internal/service/auth"
	"auth-microservice/internal/service/cache"
	"auth-microservice/internal/service/user"
	"auth-microservice/internal/transport/handler"
	"auth-microservice/internal/utils/jwt"
	"auth-microservice/internal/utils/verification"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/mailer"
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
	logger.Info(authService)

	_, s := handler.NewHandler(
		cfg,
		userService,
		authService,
	)
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.RegisterRoutes(router)

	if err := router.Run(cfg.Application.PORT); err != nil {
		logger.Panic(err, "failed to start server")
	}
}
