package http

import (
	"context"
	"github.com/aimingo/gosession/domain"
	"github.com/aimingo/gosession/service"
	"net/http"
)

// SessionMiddleware session 中间件
type SessionMiddleware struct {
	sessionService *service.SessionService
}

// NewSessionMiddleware 创建新的 session 中间件
func NewSessionMiddleware(sessionService *service.SessionService) *SessionMiddleware {
	return &SessionMiddleware{
		sessionService: sessionService,
	}
}

// MiddlewareFunc 中间件函数
func (m *SessionMiddleware) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session *domain.Session
		var err error
		var sessionID *http.Cookie
		sessionID, err = r.Cookie("session_id")
		if err == http.ErrNoCookie { // 如果 cookie 中没有 session ID，则继续下一个中间件
			session, err = m.sessionService.GetOrCreateSession(sessionID.Value)
		} else if err != nil { // 如果其他错误，则直接返回错误
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session == nil {
			session, err = m.sessionService.GetOrCreateSession(sessionID.Value)
		}

		if err != nil { // 如果获取或创建会话时出错，则直接返回错误
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 将会话对象保存到请求上下文中
		ctx := r.Context()
		ctx = context.WithValue(ctx, domain.SessionContextKey{}, session)
		r = r.WithContext(ctx)

		// 设置新的 cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			HttpOnly: true,
			Path:     "/",
		})

		// 调用下一个中间件
		next.ServeHTTP(w, r)
	})
}
