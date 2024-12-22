// auth/session.go
package auth

import (
	"net/http"
	"sync"
	"time"
	"strconv"
)

type Session struct {
	UserID       int
	LastActivity time.Time
}

type SessionStore struct {
	sessions map[string]Session
	mu       sync.Mutex
}

var store *SessionStore

// InitSessionStore инициализирует хранилище сессий
func InitSessionStore() {
	store = &SessionStore{
		sessions: make(map[string]Session),
	}
	go store.cleanupExpiredSessions()
}

// CreateSession создаёт новую сессию для пользователя
func (s *SessionStore) CreateSession(userID int) (string, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = Session{
		UserID:       userID,
		LastActivity: time.Now(),
	}
	return sessionID, nil
}

// GetSession возвращает сессию по идентификатору
func (s *SessionStore) GetSession(sessionID string) (Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, exists := s.sessions[sessionID]
	if !exists {
		return Session{}, false
	}
	// Обновление времени последней активности
	session.LastActivity = time.Now()
	s.sessions[sessionID] = session
	return session, true
}

// DeleteSession удаляет сессию по идентификатору
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

// cleanupExpiredSessions удаляет сессии, неактивные более 3 минут
func (s *SessionStore) cleanupExpiredSessions() {
	for {
		time.Sleep(time.Minute)
		s.mu.Lock()
		for id, session := range s.sessions {
			if time.Since(session.LastActivity) > 3*time.Minute {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}

// AuthMiddleware проверяет наличие действительной сессии
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Неавторизованный доступ: отсутствует сессия", http.StatusUnauthorized)
			return
		}

		session, exists := store.GetSession(cookie.Value)
		if !exists {
			http.Error(w, "Неавторизованный доступ: недействительная сессия", http.StatusUnauthorized)
			return
		}

		// Добавление UserID в заголовок запроса для дальнейшего использования
		r.Header.Set("UserID", strconv.Itoa(session.UserID))

		next.ServeHTTP(w, r)
	}
}
