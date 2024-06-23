package yandex

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/model"
	"auth-microservice/internal/service"
	"auth-microservice/internal/transport/handler/errors"
	"auth-microservice/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
)

type Handler struct {
	cfg           *config.Config
	yandexService service.YandexService
	userService   service.UserService
	jwtService    service.JwtService
}

func NewYandexHandler(
	cfg *config.Config,
	yandexService service.YandexService,
	userService service.UserService,
	jwtService service.JwtService,
) *Handler {
	return &Handler{
		cfg:           cfg,
		yandexService: yandexService,
		userService:   userService,
		jwtService:    jwtService,
	}
}

// YandexLogin godoc
// @Summary Yandex OAuth2 Login
// @Description Initiates Yandex OAuth2 login by redirecting to Yandex's consent page
// @Tags auth
// @Success 303 {string} string "Redirect to Yandex"
// @Router /api/social/yandex/login [get]
func (h Handler) YandexLogin(c *gin.Context) {
	//todo перенести в отдельный сервис под Yandex

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

	url := h.yandexService.GetYandexConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	logger.Info(url)
	c.Redirect(http.StatusSeeOther, url)
}

// YandexCallback godoc
// @Summary Yandex OAuth2 Callback
// @Description Handles the callback from Yandex OAuth2
// @Tags auth
// @Param state query string true "State"
// @Param code query string true "Authorization Code"
// @Success 200 {string} string "User data from Yandex"
// @Failure 400 {string} string "States don't Match"
// @Failure 500 {string} string "Code-Token Exchange Failed" or "User Data Fetch Failed" or "JSON Parsing Failed"
// @Router /api/social/yandex/callback [get]
func (h Handler) YandexCallback(c *gin.Context) {
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
	token, err := h.yandexService.GetYandexConfig().Exchange(context.Background(), code)
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

	resp, err := http.Get("https://login.yandex.ru/info?format=json&oauth_token=" + token.AccessToken)
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
	var user model.UserYandexInfo

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

	_, err = h.userService.CreateUserFromYandex(&model.UserCreateFromYandex{
		Id:           user.Id,
		RealName:     user.RealName,
		Avatar:       fmt.Sprintf("https://avatars.yandex.net/get-yapic/"+user.AvatarID+"/islands-200", user.AvatarID),
		IsVerified:   true,
		DefaultEmail: user.DefaultEmail,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Error: errors.Error{
				Timestamp: timestamppb.Now().String(),
				Status:    http.StatusInternalServerError,
				Error:     err.Error(),
				Message:   "error creating user",
			},
		})
		return
	}

	userInfo := model.UserInfo{
		ID:       user.Id,
		Email:    user.DefaultEmail,
		Username: user.Login,
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
		return
	}

	c.SetCookie("access_token", tokens.AccessToken, h.cfg.Jwt.ACCESS_LIFE_TIME, "/", h.cfg.Application.DOMAIN, true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, h.cfg.Jwt.REFRESH_LIFE_TIME, "/", h.cfg.Application.DOMAIN, true, true)

	c.Redirect(http.StatusSeeOther, h.cfg.Application.FRONTEND_URL)
}
