package oauth

import (
	"context"

	userpb "github.com/oj-lab/reborn/protobuf/user"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	ID        string
	Name      string
	Email     string
	AvatarURL string
}

func (u *UserInfo) ToUserPb() *userpb.User {
	return &userpb.User{
		Name:  u.Name,
		Email: u.Email,
	}
}

type Provider interface {
	GetAuthURL(state string) (string, error)
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
	Login(ctx context.Context, userInfo *UserInfo) (*userpb.User, error)
}
