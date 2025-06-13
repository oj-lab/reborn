package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"github.com/oj-lab/reborn/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
	oauthStateString = "random"
)

func LoginUser(ginCtx *gin.Context) {
	if os.Getenv("CLIENT_ID") == "" || os.Getenv("GITHUB_CLIENT_SECRET") == "" {
		ginCtx.String(http.StatusInternalServerError, "GITHUB_CLIENT_ID or GITHUB_CLIENT_SECRET is not set")
		return
	}
	url := oauthConf.AuthCodeURL(oauthStateString)
	ginCtx.Redirect(http.StatusTemporaryRedirect, url)
}

func CallbackUser(ginCtx *gin.Context) {
	state := ginCtx.Query("state")
	if state != oauthStateString {
		ginCtx.String(http.StatusBadRequest, "Invalid state")
		return
	}
	code := ginCtx.Query("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, err.Error())
		return
	}
	client := oauthConf.Client(context.Background(), token)

	// Get user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, "Failed to get user info: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var github_user utils.GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&github_user); err != nil {
		ginCtx.String(http.StatusInternalServerError, "Failed to parse user info: "+err.Error())
		return
	}

	// Get user emails
	emailsResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, "Failed to get user emails: "+err.Error())
		return
	}
	defer emailsResp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(emailsResp.Body).Decode(&emails); err != nil {
		ginCtx.String(http.StatusInternalServerError, "Failed to parse user emails: "+err.Error())
		return
	}

	// Find primary email
	for _, email := range emails {
		if email.Primary && email.Verified {
			github_user.Email = email.Email
			break
		}
	}

	// Debug log
	ginCtx.JSON(http.StatusOK, gin.H{
		"message": "GitHub User Info",
		"user":    github_user,
		"emails":  emails,
	})

	if len(emails) == 0 {
		ginCtx.String(http.StatusInternalServerError, "No primary email found")
		return
	}

	user := utils.Convert(github_user)
	req := &userpb.CreateUserRequest{
		Name:  user.Name,
		Email: emails[0].Email,
	}

	_, err = userServiceClient.CreateUser(context.Background(), req)
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	ginCtx.String(http.StatusOK, "User created successfully")
}
