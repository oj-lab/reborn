package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oj-lab/reborn/common/app"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOAuthConfig *oauth2.Config

func newGithubOAuthConfig() *oauth2.Config {
	config := app.Config()
	return &oauth2.Config{
		ClientID:     config.GetString("oauth.github.client_id"),
		ClientSecret: config.GetString("oauth.github.client_secret"),
		RedirectURL:  config.GetString("oauth.github.redirect_url"),
		Scopes:       []string{"user:email", "read:user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  github.Endpoint.AuthURL,
			TokenURL: github.Endpoint.TokenURL,
		},
	}
}

func NewGithubProvider(userServiceClient userpb.UserServiceClient) Provider {
	if githubOAuthConfig == nil {
		githubOAuthConfig = newGithubOAuthConfig()
	}

	return &GithubProvider{
		config:            githubOAuthConfig,
		userServiceClient: userServiceClient,
	}
}

type GithubProvider struct {
	config            *oauth2.Config
	userServiceClient userpb.UserServiceClient
}

func (p *GithubProvider) Login(ctx context.Context, userInfo *UserInfo) (*userpb.User, error) {
	req := &userpb.GithubLoginRequest{
		GithubId: userInfo.ID,
		Name:     userInfo.Name,
		Email:    userInfo.Email,
	}

	loginResp, err := p.userServiceClient.GithubLogin(ctx, req)
	if err != nil {
		return nil, err
	}

	return loginResp.User, nil
}

func (p *GithubProvider) GithubLoginEnabled() bool {
	if p.config.ClientID == "" || p.config.ClientSecret == "" {
		return false
	}
	return true
}

func (p *GithubProvider) GetAuthURL(state string) (string, error) {
	if !p.GithubLoginEnabled() {
		return "", fmt.Errorf("github login is not enabled")
	}
	return p.config.AuthCodeURL(state), nil
}

func (p *GithubProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *GithubProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	client := p.config.Client(context.Background(), token)

	// Get user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// In Github, the Name field might be empty. Use Login as a fallback.
	if githubUser.Name == "" {
		githubUser.Name = githubUser.Login
	}

	// If email is not available in user info, get from emails endpoint
	if githubUser.Email == "" {
		emailsResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, fmt.Errorf("failed to get user emails: %w", err)
		}
		defer emailsResp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := json.NewDecoder(emailsResp.Body).Decode(&emails); err != nil {
			return nil, fmt.Errorf("failed to parse user emails: %w", err)
		}

		for _, email := range emails {
			if email.Primary && email.Verified {
				githubUser.Email = email.Email
				break
			}
		}
	}

	if githubUser.Email == "" {
		return nil, fmt.Errorf("github user email is not available")
	}

	return &UserInfo{
		ID:        fmt.Sprintf("%d", githubUser.ID),
		Name:      githubUser.Name,
		Email:     githubUser.Email,
		AvatarURL: githubUser.AvatarURL,
	}, nil
}
