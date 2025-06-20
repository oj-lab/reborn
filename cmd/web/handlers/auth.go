package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/common/oauth"
	"github.com/oj-lab/reborn/common/redis_client"
	"github.com/oj-lab/reborn/common/session"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"github.com/redis/go-redis/v9"
)

const (
	sessionTTL    = 24 * time.Hour
	oauthStateTTL = 10 * time.Minute
)

var (
	sessionManager = session.NewManager()
	rdb            = redis_client.GetRDB()
)

func getOauthStateKey(state string) string {
	return fmt.Sprintf("oauth:state:%s", state)
}

func Login(ctx echo.Context) error {
	providerName := ctx.Param("provider")
	provider, err := oauth.GetProvider(providerName)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	oauthState := uuid.New().String()
	stateKey := getOauthStateKey(oauthState)
	if err := rdb.Set(ctx.Request().Context(), stateKey, "true", oauthStateTTL).Err(); err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to save state")
	}

	url, err := provider.GetAuthURL(oauthState)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func Callback(ctx echo.Context) error {
	providerName := ctx.Param("provider")
	provider, err := oauth.GetProvider(providerName)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	state := ctx.QueryParam("state")
	stateKey := getOauthStateKey(state)
	err = rdb.Get(ctx.Request().Context(), stateKey).Err()
	if err == redis.Nil {
		return ctx.String(http.StatusBadRequest, "Invalid or expired state")
	} else if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to verify state")
	}
	rdb.Del(ctx.Request().Context(), stateKey)

	code := ctx.QueryParam("code")
	token, err := provider.Exchange(context.Background(), code)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
	}

	userInfo, err := provider.GetUserInfo(token)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to get user info: "+err.Error())
	}

	if providerName != "github" {
		return ctx.String(http.StatusBadRequest, "Unsupported provider")
	}

	req := &userpb.GithubLoginRequest{
		GithubId: userInfo.ID,
		Name:     userInfo.Name,
		Email:    userInfo.Email,
	}

	loginResp, err := userServiceClient.GithubLogin(context.Background(), req)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to login: "+err.Error())
	}

	sessionID, err := sessionManager.Create(context.Background(), uint(loginResp.User.Id), sessionTTL)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create session: "+err.Error())
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(sessionTTL),
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusOK, map[string]string{"message": "login success"})
}

type PasswordLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginWithPassword(c echo.Context) error {
	req := new(PasswordLoginRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	grpcReq := &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := userServiceClient.Login(context.Background(), grpcReq)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to login: "+err.Error())
	}

	sessionID, err := sessionManager.Create(context.Background(), uint(resp.User.Id), sessionTTL)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create session: "+err.Error())
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(sessionTTL),
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "login success"})
}

type PasswordRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterWithPassword(c echo.Context) error {
	req := new(PasswordRegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	grpcReq := &userpb.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: &req.Password,
	}

	_, err := userServiceClient.CreateUser(context.Background(), grpcReq)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to register: "+err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

type SetPasswordRequest struct {
	Password string `json:"password"`
}

func SetPassword(c echo.Context) error {
	req := new(SetPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	userID, ok := c.Get("userID").(uint64)
	if !ok || userID == 0 {
		return c.String(http.StatusUnauthorized, "Invalid user ID from token")
	}

	grpcReq := &userpb.SetPasswordRequest{
		UserId:   userID,
		Password: req.Password,
	}

	_, err := userServiceClient.SetPassword(context.Background(), grpcReq)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to set password: "+err.Error())
	}

	return c.NoContent(http.StatusOK)
}
