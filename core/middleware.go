package main

import (
	"context"
	"log"
	"net/http"
)

// contextKey is a custom type for context keys
type contextKey string

// UserContextKey is the key for storing user in request context
const UserContextKey contextKey = "user"

// AuthMiddleware validates JWT and adds user to request context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT token from cookie or header
		tokenString, err := jwtService.ExtractToken(r)
		if err != nil {
			log.Printf("Failed to extract token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Load user from database
		user, err := GetUserByID(claims.UserID.String())
		if err != nil {
			log.Printf("Failed to load user: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r = r.WithContext(ctx)

		// Call next handler
		next(w, r)
	}
}

// GetUserFromContext retrieves the user from request context
func GetUserFromContext(r *http.Request) (*User, error) {
	user, ok := r.Context().Value(UserContextKey).(*User)
	if !ok {
		return nil, http.ErrNoCookie
	}
	return user, nil
}
