package oauth

import (
	"fmt"

	userpb "github.com/oj-lab/reborn/protobuf/user"
)

func GetProvider(name string, userServiceClient userpb.UserServiceClient) (Provider, error) {
	switch name {
	case "github":
		return NewGithubProvider(userServiceClient), nil
	default:
		return nil, fmt.Errorf("provider %s not supported", name)
	}
}
