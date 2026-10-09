package sessionmanager

import (
	"time"
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

type sessionManager struct{
	ttl    time.Duration
	maxAge int
	isProd bool
}

func New(ttl time.Duration) *sessionManager {
	return &sessionManager{
		ttl:    ttl,
		maxAge: int(ttl.Seconds()),
	}
}

func (m *sessionManager) Generate() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (m *sessionManager) Cookie(id string) *http.Cookie {
	return &http.Cookie{
		Name: "session",
		Value: id,
		MaxAge: m.maxAge,
		SameSite: http.SameSiteLaxMode,
		Secure: m.isProd,
		HttpOnly: true,
	}
}

func (m *sessionManager) TTL() time.Duration {
	return m.ttl
}
