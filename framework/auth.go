package framework

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"flux/config"
)

// TokenPurpose distinguishes the different signed tokens Flux issues.
type TokenPurpose string

const (
	PurposeAccess        TokenPurpose = "access"
	PurposeEmailVerify   TokenPurpose = "email_verify"
	PurposePasswordReset TokenPurpose = "password_reset"
)

// Claims are the JWT claims carried by every Flux token.
type Claims struct {
	UserID  uint         `json:"uid"`
	Email   string       `json:"email"`
	Purpose TokenPurpose `json:"purpose"`
	jwt.RegisteredClaims
}

// AuthManager issues and verifies JWTs (signed with APP_KEY) and manages the
// session cookie that carries the access token for browser navigation.
type AuthManager struct {
	cfg *config.Config
}

func NewAuthManager(cfg *config.Config) *AuthManager {
	return &AuthManager{cfg: cfg}
}

// SessionCookieName is the name of the httpOnly cookie holding the JWT.
func (m *AuthManager) SessionCookieName() string { return m.cfg.Auth.SessionCookie }

// SessionTTL is how long issued access tokens (and the cookie) live.
func (m *AuthManager) SessionTTL() time.Duration {
	return time.Duration(m.cfg.Auth.TokenTTLMinutes) * time.Minute
}

// IssueAccessToken creates the JWT used for both API auth (Bearer header)
// and browser sessions (httpOnly cookie).
func (m *AuthManager) IssueAccessToken(userID uint, email string) (string, error) {
	return m.issue(userID, email, PurposeAccess, m.SessionTTL())
}

// IssueEmailVerifyToken creates a short-lived signed token for email
// verification links.
func (m *AuthManager) IssueEmailVerifyToken(userID uint, email string) (string, error) {
	return m.issue(userID, email, PurposeEmailVerify, 24*time.Hour)
}

// IssuePasswordResetToken creates a short-lived signed token for password
// reset links.
func (m *AuthManager) IssuePasswordResetToken(userID uint, email string) (string, error) {
	return m.issue(userID, email, PurposePasswordReset, time.Hour)
}

func (m *AuthManager) issue(userID uint, email string, purpose TokenPurpose, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:  userID,
		Email:   email,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.App.Name,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.App.Key))
}

// ParseToken verifies a token's signature, expiry and purpose.
func (m *AuthManager) ParseToken(tokenString string, purpose TokenPurpose) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(m.cfg.App.Key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Purpose != purpose {
		return nil, fmt.Errorf("token purpose mismatch: got %q, want %q", claims.Purpose, purpose)
	}
	return claims, nil
}

// TokenFromRequest extracts the access token from the Authorization header
// (Bearer) or, failing that, the session cookie.
func (m *AuthManager) TokenFromRequest(r *Request) string {
	const bearerPrefix = "Bearer "
	if auth := r.Header("Authorization"); len(auth) > len(bearerPrefix) && auth[:len(bearerPrefix)] == bearerPrefix {
		return auth[len(bearerPrefix):]
	}
	return r.Cookie(m.cfg.Auth.SessionCookie)
}

// StartSession sets the httpOnly session cookie carrying the JWT.
func (m *AuthManager) StartSession(r *Request, token string) {
	m.setSessionCookie(r, token, int(m.SessionTTL().Seconds()))
}

// EndSession clears the session cookie.
func (m *AuthManager) EndSession(r *Request) {
	m.setSessionCookie(r, "", -1)
}

func (m *AuthManager) setSessionCookie(r *Request, value string, maxAge int) {
	r.gin.SetCookie(m.cfg.Auth.SessionCookie, value, maxAge, "/", "", secureCookies(m.cfg), true)
}

// secureCookies reports whether cookies should carry the Secure flag —
// shared by the session and flash cookies.
func secureCookies(cfg *config.Config) bool { return cfg.App.Env == "production" }

// HashPassword hashes a plaintext password with bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword reports whether a plaintext password matches a bcrypt hash.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
