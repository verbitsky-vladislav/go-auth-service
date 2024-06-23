package yandex

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/service"
	"golang.org/x/oauth2"
)

type yandexService struct {
	cfg               *config.Config
	yandexOauthConfig *oauth2.Config
	userService       service.UserService
}

func NewYandexService(
	cfg *config.Config,
	userService service.UserService,
) service.YandexService {
	return &yandexService{
		cfg: cfg,
		yandexOauthConfig: &oauth2.Config{
			RedirectURL:  cfg.Yandex.REDIRECT_URL,
			ClientID:     cfg.Yandex.CLIENT_ID,
			ClientSecret: cfg.Yandex.CLIENT_SECRET,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.yandex.com/authorize",
				TokenURL: "https://oauth.yandex.com/token",
			},
			Scopes: []string{},
		},
		userService: userService,
	}
}

func (s *yandexService) GetYandexConfig() *oauth2.Config {
	return s.yandexOauthConfig
}
