package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"github.com/oj-lab/reborn/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
	oauthStateString = "random"
)

func LoginUser(ctx echo.Context) error {
	if os.Getenv("CLIENT_ID") == "" || os.Getenv("GITHUB_CLIENT_SECRET") == "" {
		return ctx.String(http.StatusInternalServerError, "GitHub OAuth credentials are not set")
	}
	url := oauthConf.AuthCodeURL(oauthStateString)
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func CallbackUser(ctx echo.Context) error {
	state := ctx.QueryParam("state")
	if state != oauthStateString {
		return ctx.String(http.StatusBadRequest, "Invalid state")
	}
	code := ctx.QueryParam("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	client := oauthConf.Client(context.Background(), token)

	// Get user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to get user info: "+err.Error())
	}
	defer resp.Body.Close()

	var github_user utils.GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&github_user); err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to parse user info: "+err.Error())
	}

	// Get user emails
	emailsResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to get user emails: "+err.Error())
	}
	defer emailsResp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(emailsResp.Body).Decode(&emails); err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to parse user emails: "+err.Error())
	}

	// Find primary email
	for _, email := range emails {
		if email.Primary && email.Verified {
			github_user.Email = email.Email
			break
		}
	}

	// Debug log
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "GitHub User Info",
		"user":    github_user,
		"emails":  emails,
	})

	if len(emails) == 0 {
		return ctx.String(http.StatusInternalServerError, "No primary email found")
	}

	user := utils.Convert(github_user)
	req := &userpb.CreateUserRequest{
		Name:  user.Name,
		Email: emails[0].Email,
	}

	_, err = userServiceClient.CreateUser(context.Background(), req)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create user: "+err.Error())
	}
	return ctx.String(http.StatusOK, "User created successfully")
}
