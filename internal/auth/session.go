// Package auth implements the single shared-password gate. It verifies the
// password and issues/validates a signed, stateless session cookie
// (HMAC-SHA256 over an expiry timestamp), so no session store is needed.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	cookieName = "session"
	sessionTTL = 7 * 24 * time.Hour
)

// Manager verifies the shared password and manages session cookies.
type Manager struct {
	secret   []byte
	password string
	secure   bool
}

func NewManager(secret, password string, secure bool) *Manager {
	return &Manager{secret: []byte(secret), password: password, secure: secure}
}

// CheckPassword reports whether pw matches the configured shared password,
// using a constant-time comparison to avoid timing leaks.
func (m *Manager) CheckPassword(pw string) bool {
	return subtle.ConstantTimeCompare([]byte(pw), []byte(m.password)) == 1
}

// Issue writes a fresh signed session cookie to the response.
func (m *Manager) Issue(w http.ResponseWriter) {
	expiry := time.Now().Add(sessionTTL)
	http.SetCookie(w, m.cookie(m.sign(expiry), expiry))
}

// Clear expires the session cookie (logout).
func (m *Manager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, m.cookie("", time.Unix(0, 0)))
}

// Valid reports whether the request carries a valid, unexpired session cookie.
func (m *Manager) Valid(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	return m.verify(c.Value)
}

func (m *Manager) cookie(value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	}
}

// sign produces "base64(payload).base64(hmac)" where payload is the expiry.
func (m *Manager) sign(expiry time.Time) string {
	payload := strconv.FormatInt(expiry.Unix(), 10)
	return encode(payload) + "." + encode(string(m.mac(payload)))
}

func (m *Manager) verify(token string) bool {
	encPayload, encSig, ok := strings.Cut(token, ".")
	if !ok {
		return false
	}
	payload, err := decode(encPayload)
	if err != nil {
		return false
	}
	sig, err := decode(encSig)
	if err != nil {
		return false
	}
	if !hmac.Equal([]byte(sig), m.mac(payload)) {
		return false
	}
	expiryUnix, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return false
	}
	return time.Now().Before(time.Unix(expiryUnix, 0))
}

func (m *Manager) mac(payload string) []byte {
	h := hmac.New(sha256.New, m.secret)
	h.Write([]byte(payload))
	return h.Sum(nil)
}

func encode(s string) string { return base64.RawURLEncoding.EncodeToString([]byte(s)) }

func decode(s string) (string, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	return string(b), err
}
