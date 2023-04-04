package infrastructure

import (
	"container/list"
	"github.com/aimingo/gosession/domain"
	"sync"
	"time"
)

// SessionRepositoryImpl 会话仓储实现
type SessionRepositoryImpl struct {
	mu          sync.RWMutex
	sessionMap  map[string]*list.Element // 会话存储的 map
	sessionList *list.List               // 会话列表
}

// NewSessionRepositoryImpl 创建新的会话仓储实现
func NewSessionRepositoryImpl() *SessionRepositoryImpl {
	return &SessionRepositoryImpl{
		sessionMap:  make(map[string]*list.Element),
		sessionList: list.New(),
	}
}

// Get 根据会话 ID 获取会话实体对象
func (repo *SessionRepositoryImpl) Get(id string) (*domain.Session, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if elem, ok := repo.sessionMap[id]; ok {
		session := elem.Value.(*domain.Session)
		session.ActiveAt = time.Now()      // 刷新会话的活跃时间
		repo.sessionList.MoveToFront(elem) // 将会话移动到列表头部，表示最近使用
		return session, nil
	}

	return nil, domain.ErrSessionNotFound
}

// Set 保存会话实体对象
func (repo *SessionRepositoryImpl) Set(session *domain.Session) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	// 如果超过最大会话数，则使用 LRU 淘汰
	if repo.sessionList.Len() >= 100000 {
		repo.evictLRU()
	}

	// 如果会话已存在，则更新会话
	if elem, ok := repo.sessionMap[session.ID]; ok {
		elem.Value = session
		repo.sessionList.MoveToFront(elem)
	} else { // 否则添加新的会话
		elem := repo.sessionList.PushFront(session)
		repo.sessionMap[session.ID] = elem
	}

	return nil
}

// Delete 删除会话
func (repo *SessionRepositoryImpl) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if elem, ok := repo.sessionMap[id]; ok {
		delete(repo.sessionMap, id)
		repo.sessionList.Remove(elem)
		return nil
	}

	return domain.ErrSessionNotFound
}

// Len 获取当前会话数
func (repo *SessionRepositoryImpl) Len() int {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	return repo.sessionList.Len()
}

// evictLRU 使用 LRU 算法淘汰最久未使用的会话
func (repo *SessionRepositoryImpl) evictLRU() {
	elem := repo.sessionList.Back()
	if elem != nil {
		session := elem.Value.(*domain.Session)
		repo.Delete(session.ID)
	}
}

// PeekOldest 获取最久未使用的会话
func (repo *SessionRepositoryImpl) PeekOldest() (*domain.Session, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	elem := repo.sessionList.Back()
	if elem != nil {
		session := elem.Value.(*domain.Session)
		return session, nil
	}

	return nil, domain.ErrSessionNotFound
}
