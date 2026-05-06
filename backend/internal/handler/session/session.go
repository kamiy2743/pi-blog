package session

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"blog/internal/app"
	"blog/internal/handler/handlererror"
	"blog/internal/handler/middleware"

	"github.com/alexedwards/scs/v2"
)

const (
	cookieName          = "blog_session"
	lifetime            = 5 * time.Minute
	idleTimeout         = 3 * time.Minute
	oldInputKey         = "oldInput"
	validationErrorsKey = "validationErrors"
	flashSuccessKey     = "flash.success"
	flashErrorKey       = "flash.error"
)

type SessionManager struct {
	manager *scs.SessionManager
}

type SessionPayload struct {
	OldInput        map[string]string
	ValidationError *handlererror.ValidationError
	Flash           *Flash
}

type contextKey struct{}

func NewSessionManager(appEnv app.AppEnv) *SessionManager {
	manager := scs.New()
	manager.Lifetime = lifetime
	manager.IdleTimeout = idleTimeout
	manager.Cookie.Name = cookieName
	manager.Cookie.HttpOnly = true
	manager.Cookie.Secure = appEnv == app.AppEnvPrd
	manager.Cookie.SameSite = http.SameSiteLaxMode

	return &SessionManager{
		manager: manager,
	}
}

func (m *SessionManager) Middleware() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return m.manager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextKey{}, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		}))
	}
}

func SessionManagerFromContext(ctx context.Context) (*SessionManager, bool) {
	manager, ok := ctx.Value(contextKey{}).(*SessionManager)
	return manager, ok
}

func (m *SessionManager) SaveOldInput(r *http.Request, oldInput map[string]string) {
	if len(oldInput) == 0 {
		return
	}

	oldInputJSON, err := json.Marshal(oldInput)
	if err != nil {
		return
	}
	m.manager.Put(r.Context(), oldInputKey, string(oldInputJSON))
}

func (m *SessionManager) SaveValidationError(r *http.Request, validationError *handlererror.ValidationError) {
	if validationError == nil || validationError.IsEmpty() {
		return
	}

	validationErrorsJSON, err := json.Marshal(validationError.Messages)
	if err != nil {
		return
	}
	m.manager.Put(r.Context(), validationErrorsKey, string(validationErrorsJSON))
}

func (m *SessionManager) SaveFlash(r *http.Request, flash *Flash) {
	if flash == nil || flash.IsEmpty() {
		return
	}

	if flash.Success != "" {
		m.manager.Put(r.Context(), flashSuccessKey, flash.Success)
	}
	if flash.Error != "" {
		m.manager.Put(r.Context(), flashErrorKey, flash.Error)
	}
}

func (m *SessionManager) PopSessionPayload(r *http.Request) SessionPayload {
	return m.sessionPayload(r, true)
}

func (m *SessionManager) GetSessionPayload(r *http.Request) SessionPayload {
	return m.sessionPayload(r, false)
}

func (m *SessionManager) sessionPayload(r *http.Request, pop bool) SessionPayload {
	oldInput := map[string]string{}
	oldInputJSON := m.manager.GetString(r.Context(), oldInputKey)
	if pop {
		oldInputJSON = m.manager.PopString(r.Context(), oldInputKey)
	}
	if oldInputJSON != "" {
		json.Unmarshal([]byte(oldInputJSON), &oldInput)
	}

	var validationError *handlererror.ValidationError
	validationErrorsJSON := m.manager.GetString(r.Context(), validationErrorsKey)
	if pop {
		validationErrorsJSON = m.manager.PopString(r.Context(), validationErrorsKey)
	}
	if validationErrorsJSON != "" {
		validationErrors := map[string]string{}
		json.Unmarshal([]byte(validationErrorsJSON), &validationErrors)
		validationError = &handlererror.ValidationError{Messages: validationErrors}
	}

	flash := &Flash{
		Success: m.manager.GetString(r.Context(), flashSuccessKey),
		Error:   m.manager.GetString(r.Context(), flashErrorKey),
	}
	if pop {
		flash = &Flash{
			Success: m.manager.PopString(r.Context(), flashSuccessKey),
			Error:   m.manager.PopString(r.Context(), flashErrorKey),
		}
	}
	if flash.IsEmpty() {
		flash = nil
	}

	return SessionPayload{
		OldInput:        oldInput,
		ValidationError: validationError,
		Flash:           flash,
	}
}
