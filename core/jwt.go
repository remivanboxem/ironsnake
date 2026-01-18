package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	config *JWTConfig
}

// NewJWTService creates a new JWT service instance
func NewJWTService(config *JWTConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateToken creates a new JWT token for a user
func (s *JWTService) GenerateToken(user *User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.config.ExpirationHours) * time.Hour)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ironsnake",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ExtractTokenFromCookie extracts the JWT token from the request cookie
func (s *JWTService) ExtractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", fmt.Errorf("auth token cookie not found")
	}

	if cookie.Value == "" {
		return "", fmt.Errorf("auth token cookie is empty")
	}

	return cookie.Value, nil
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header
func (s *JWTService) ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	// Expected format: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

// ExtractToken tries to extract token from cookie first, then from header
func (s *JWTService) ExtractToken(r *http.Request) (string, error) {
	// Try cookie first
	token, err := s.ExtractTokenFromCookie(r)
	if err == nil {
		return token, nil
	}

	// Fallback to authorization header
	token, err = s.ExtractTokenFromHeader(r)
	if err != nil {
		return "", fmt.Errorf("no valid token found in cookie or header")
	}

	return token, nil
}

// SetTokenCookie sets the JWT token as an HTTP-only cookie
func (s *JWTService) SetTokenCookie(w http.ResponseWriter, token string, secure bool) {
	maxAge := s.config.ExpirationHours * 3600 // Convert hours to seconds

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearTokenCookie removes the JWT token cookie
func (s *JWTService) ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
