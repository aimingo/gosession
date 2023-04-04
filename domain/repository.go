package domain

// SessionRepository 会话仓储接口定义
type SessionRepository interface {
	Get(id string) (*Session, error)
	Set(session *Session) error
	Delete(id string) error
	Len() int
	PeekOldest() (*Session, error)
}
