package google

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleService struct {
	cfg               *config.Config
	googleOauthConfig *oauth2.Config
	userService       service.UserService
}

func NewGoogleService(
	cfg *config.Config,
	userService service.UserService,
) service.GoogleService {
	return &googleService{
		cfg: cfg,
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  cfg.Google.REDIRECT_URL,
			ClientID:     cfg.Google.CLIENT_ID,
			ClientSecret: cfg.Google.CLIENT_SECRET,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/photoslibrary.readonly",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		userService: userService,
	}
}

func (s *googleService) GetGoogleConfig() *oauth2.Config {
	return s.googleOauthConfig
}
