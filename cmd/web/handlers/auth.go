package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/common/oauth"
	userpb "github.com/oj-lab/reborn/protobuf/user"
)

const (
	oauthStateString = "random" // TODO: Should be a random string
)

func Login(ctx echo.Context) error {
	providerName := ctx.Param("provider")
	provider, err := oauth.GetProvider(providerName)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	url, err := provider.GetAuthURL(oauthStateString)
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
	if state != oauthStateString {
		return ctx.String(http.StatusBadRequest, "Invalid state")
	}

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

	// TODO: Login user and set session
	return ctx.JSON(http.StatusOK, loginResp)
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

	return c.JSON(http.StatusOK, resp)
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
