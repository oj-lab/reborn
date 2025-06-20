package oauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// State represents the state parameter used in the OAuth flow.
// It includes a CSRF token and the provider name.
type State struct {
	CSRFToken string `json:"csrf_token"`
	Provider  string `json:"provider"`
}

// NewState creates a new State object with a random CSRF token.
func NewState(provider string) *State {
	return &State{
		CSRFToken: uuid.NewString(),
		Provider:  provider,
	}
}

// Encode serializes the State to a base64-encoded JSON string.
func (s *State) Encode() (string, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(jsonData), nil
}

// DecodeState deserializes a State object from a base64-encoded JSON string.
func DecodeState(encodedState string) (*State, error) {
	jsonData, err := base64.URLEncoding.DecodeString(encodedState)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state: %w", err)
	}

	var state State
	if err := json.Unmarshal(jsonData, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}
	return &state, nil
}
