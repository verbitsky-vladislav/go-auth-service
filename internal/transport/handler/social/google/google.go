package google

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/model"
	"auth-microservice/internal/service"
	"auth-microservice/internal/transport/handler/errors"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
)

type Handler struct {
	cfg           *config.Config
	googleService service.GoogleService
	userService   service.UserService
	jwtService    service.JwtService
}

func NewGoogleHandler(
	cfg *config.Config,
	googleService service.GoogleService,
	userService service.UserService,
	jwtService service.JwtService,
) *Handler {
	return &Handler{
		cfg:           cfg,
		googleService: googleService,
		userService:   userService,
		jwtService:    jwtService,
	}
}

// GoogleLogin godoc
// @Summary Google OAuth2 Login
// @Description Initiates Google OAuth2 login by redirecting to Google's consent page
// @Tags auth
// @Success 303 {string} string "Redirect to Google"
// @Router /api/google/login [get]
func (h Handler) GoogleLogin(c *gin.Context) {
	//todo перенести в отдельный сервис под гугл

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "failed to generate random string",
			},
		})
		return
	}
	state := base64.URLEncoding.EncodeToString(b)

	c.SetCookie("oauthstate", state, 3600, "/", "", false, true)

	url := h.googleService.GetGoogleConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusSeeOther, url)
}

// GoogleCallback godoc
// @Summary Google OAuth2 Callback
// @Description Handles the callback from Google OAuth2
// @Tags auth
// @Param state query string true "State"
// @Param code query string true "Authorization Code"
// @Success 200 {string} string "User data from Google"
// @Failure 400 {string} string "States don't Match!!"
// @Failure 500 {string} string "Code-Token Exchange Failed" or "User Data Fetch Failed" or "JSON Parsing Failed"
// @Router /api/google/callback [get]
func (h Handler) GoogleCallback(c *gin.Context) {
	//todo перенести в отдельный сервис под гугл

	state := c.Query("state")
	oauthState, err := c.Cookie("oauthstate")
	if err != nil || state != oauthState {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "invalid oauth state",
			},
		})
		return
	}

	code := c.Query("code")
	token, err := h.googleService.GetGoogleConfig().Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "invalid oauth code",
			},
		})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "error fetching user info",
			},
		})
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "error reading user info",
			},
		})
		return
	}

	var user model.UserGoogleInfo

	if err := json.Unmarshal(userData, &user); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "error parsing user info",
			},
		})
		return
	}

	id, err := h.userService.CreateUserFromGoogle(&model.UserCreateFromGoogle{
		Id:         user.Id,
		Email:      user.Email,
		Logo:       user.Picture,
		Username:   user.Name,
		IsVerified: true,
	})
	if err != nil {
		if err.Error() == "user already exists" {

		}
	}

	userInfo := model.UserInfo{
		ID:       id,
		Email:    user.Email,
		Username: user.Name,
	}

	tokens, err := h.jwtService.GenerateTokens(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusInternalServerError,
				Error:     err.Error(),
				Message:   "error generating tokens",
			},
		})
	}

	c.SetCookie("access_token", tokens.AccessToken, h.cfg.Jwt.ACCESS_LIFE_TIME, "/", h.cfg.Application.DOMAIN, true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, h.cfg.Jwt.REFRESH_LIFE_TIME, "/", h.cfg.Application.DOMAIN, true, true)

	c.Redirect(http.StatusSeeOther, h.cfg.Application.FRONTEND_URL)
}
