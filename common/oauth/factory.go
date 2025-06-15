package oauth

import "fmt"

func GetProvider(name string) (Provider, error) {
	switch name {
	case "github":
		return NewGitHubProvider(), nil
	default:
		return nil, fmt.Errorf("provider %s not supported", name)
	}
}
