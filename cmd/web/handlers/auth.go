package handlers

import (
	"context"
	"fmt"
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

	fmt.Println(userInfo)
	// TODO: Check if user exists, if not create user
	req := &userpb.CreateUserRequest{
		Name:  userInfo.Name,
		Email: userInfo.Email,
	}

	_, err = userServiceClient.CreateUser(context.Background(), req)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create user: "+err.Error())
	}

	// TODO: Login user and set session
	return ctx.JSON(http.StatusOK, userInfo)
}
