package utils

import (
	"github.com/google/uuid"
	userpb "github.com/oj-lab/reborn/protobuf/user"
)

type GithubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func Convert(user GithubUser) userpb.User {
	return userpb.User{
		Id:    uint64(uuid.New().ID()),
		Name:  user.Name,
		Email: user.Email,
	}
}
