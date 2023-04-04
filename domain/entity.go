package domain

import "time"

// Session 会话实体对象
type Session struct {
	ID        string
	Data      map[string]interface{}
	CreatedAt time.Time
	ActiveAt  time.Time
}

// NewSession 创建新的 Session 实体对象
func NewSession(id string) *Session {
	return &Session{
		ID:        id,
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		ActiveAt:  time.Now(),
	}
}

// UpdateLastActiveTime 更新会话的最后活跃时间
func (s *Session) UpdateLastActiveTime() {
	s.ActiveAt = time.Now()
}

type SessionContextKey struct {
}
