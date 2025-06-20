package session

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/oj-lab/reborn/common/redis_client"
	"github.com/redis/go-redis/v9"
)

const DefaultSessionTTL = 24 * time.Hour

type Manager struct {
	rdb redis.UniversalClient
	ttl time.Duration
}

func NewManager() *Manager {
	return &Manager{
		rdb: redis_client.GetRDB(),
		ttl: DefaultSessionTTL,
	}
}

type Session struct {
	ID     string
	UserID uint
	Data   map[string]any
}

func (m *Manager) Create(ctx context.Context, userID uint, ttl time.Duration) (string, error) {
	sessionID := uuid.New().String()
	session := &Session{
		ID:     sessionID,
		UserID: userID,
		Data:   make(map[string]any),
	}

	if err := m.rdb.JSONSet(ctx, sessionID, ".", session).Err(); err != nil {
		return "", err
	}

	if err := m.rdb.Expire(ctx, sessionID, ttl).Err(); err != nil {
		// Clean up the key if we fail to set TTL, to avoid session staying forever.
		if delErr := m.rdb.Del(ctx, sessionID).Err(); delErr != nil {
			slog.WarnContext(ctx, "Failed to delete session key after failing to set TTL", "sessionID", sessionID, "error", delErr)
		}
		return "", err
	}

	return sessionID, nil
}

func (m *Manager) Get(ctx context.Context, sessionID string) (*Session, error) {
	val, err := m.rdb.JSONGet(ctx, sessionID, ".").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Session not found
		}
		return nil, err
	}

	err = m.rdb.Expire(ctx, sessionID, m.ttl).Err()
	if err != nil {
		slog.WarnContext(ctx, "Failed to extend session TTL", "sessionID", sessionID, "error", err)
	}

	var sessions []Session
	if err := json.Unmarshal([]byte(val), &sessions); err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}

func (m *Manager) Delete(ctx context.Context, sessionID string) error {
	return m.rdb.Del(ctx, sessionID).Err()
}
