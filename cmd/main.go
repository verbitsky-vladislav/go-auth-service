package main

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/repository/postgres"
	"auth-microservice/internal/service/cache"
	"auth-microservice/internal/service/user"
	"auth-microservice/internal/transport/handler"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/mailer"
	"github.com/gin-gonic/gin"
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

	redisClient, err := cache.NewRedisClient(
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
	logger.Info(userService)
	logger.Info(mailerService)

	_, s := handler.NewHandler()
	router := gin.Default()

	s.RegisterRoutes(router)

	if err := router.Run(cfg.Application.PORT); err != nil {
		logger.Panic(err, "failed to start server")
	}
}
