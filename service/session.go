package service

import (
	"github.com/aimingo/gosession/domain"
	"github.com/aimingo/gosession/util"
	"time"
)

// SessionService 会话管理服务
type SessionService struct {
	repo domain.SessionRepository
}

// NewSessionService 创建新的会话管理服务
func NewSessionService(repo domain.SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

// GetOrCreateSession 根据 cookie 中的 session ID 获取会话实体对象，如果不存在则创建新的会话
func (s *SessionService) GetOrCreateSession(sessionID string) (*domain.Session, error) {
	if sessionID == "" { // 如果 session ID 为空，则创建新的会话
		sessionID = util.GenerateSessionID()
	}

	session, err := s.repo.Get(sessionID)
	if err != nil {
		if err == domain.ErrSessionNotFound { // 如果会话不存在，则创建新的会话
			session = domain.NewSession(sessionID)
			if err := s.repo.Set(session); err != nil {
				return nil, err
			}
		} else { // 如果其他错误，则直接返回错误
			return nil, err
		}
	}

	return session, nil
}

// DeleteSession 删除会话
func (s *SessionService) DeleteSession(sessionID string) error {
	return s.repo.Delete(sessionID)
}

// CleanupExpiredSessions 清理过期会话
func (s *SessionService) CleanupExpiredSessions() {
	for {
		time.Sleep(30 * time.Minute) // 每隔 30 分钟清理一次过期会话
		s.cleanup()
	}
}

// cleanup 清理过期会话
func (s *SessionService) cleanup() {
	now := time.Now()
	for {
		session, err := s.repo.PeekOldest()
		if err != nil { // 如果会话列表为空，则退出循环
			break
		}

		// 如果会话已过期，则删除会话并继续循环
		if session.ActiveAt.Add(1 * time.Hour).Before(now) {
			s.repo.Delete(session.ID)
		} else { // 否则退出循环
			break
		}
	}
}
