package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/cmd/web/middleware"
	"github.com/oj-lab/reborn/common/oauth"
	"github.com/oj-lab/reborn/common/redis_client"
	"github.com/oj-lab/reborn/common/session"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"github.com/redis/go-redis/v9"
)

const (
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
	providerName := ctx.QueryParam("provider")
	if providerName == "" {
		return ctx.String(http.StatusBadRequest, "provider is required")
	}

	provider, err := oauth.GetProvider(providerName, userServiceClient)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	state := oauth.NewState(providerName)
	encodedState, err := state.Encode()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to create state")
	}

	stateKey := getOauthStateKey(state.CSRFToken)
	if err := rdb.Set(ctx.Request().Context(), stateKey, "true", oauthStateTTL).Err(); err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to save state")
	}

	url, err := provider.GetAuthURL(encodedState)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func Callback(ctx echo.Context) error {
	encodedState := ctx.QueryParam("state")
	state, err := oauth.DecodeState(encodedState)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid state format")
	}

	provider, err := oauth.GetProvider(state.Provider, userServiceClient)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	stateKey := getOauthStateKey(state.CSRFToken)
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

	user, err := provider.Login(context.Background(), userInfo)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to login: "+err.Error())
	}

	sessionID, err := sessionManager.Create(context.Background(), uint(user.Id), session.DefaultSessionTTL)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create session: "+err.Error())
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(session.DefaultSessionTTL),
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

	sessionID, err := sessionManager.Create(context.Background(), uint(resp.User.Id), session.DefaultSessionTTL)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create session: "+err.Error())
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(session.DefaultSessionTTL),
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

	session, ok := c.Get(middleware.ContextKeyUserSession).(*session.Session)
	if !ok || session == nil {
		return c.String(http.StatusUnauthorized, "Invalid user ID from token")
	}

	grpcReq := &userpb.SetPasswordRequest{
		UserId:   uint64(session.UserID),
		Password: req.Password,
	}

	_, err := userServiceClient.SetPassword(context.Background(), grpcReq)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to set password: "+err.Error())
	}

	return c.NoContent(http.StatusOK)
}
