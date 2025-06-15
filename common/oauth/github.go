package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider() Provider {
	return &GithubProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/auth/github/callback",
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		},
	}
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
		ID:        fmt.Sprintf("github_%d", githubUser.ID),
		Name:      githubUser.Name,
		Email:     githubUser.Email,
		AvatarURL: githubUser.AvatarURL,
	}, nil
}
