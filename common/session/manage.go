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

	// TODO: Use Redis JSON support if available
	val, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	if err := m.rdb.Set(ctx, sessionID, val, ttl).Err(); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (m *Manager) Get(ctx context.Context, sessionID string) (*Session, error) {
	val, err := m.rdb.Get(ctx, sessionID).Bytes()
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

	var session Session
	if err := json.Unmarshal(val, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (m *Manager) Delete(ctx context.Context, sessionID string) error {
	return m.rdb.Del(ctx, sessionID).Err()
}
