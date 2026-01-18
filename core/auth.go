package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	ldapService *LDAPService
	jwtService  *JWTService
)

// InitAuthServices initializes the authentication services
func InitAuthServices(config *Config) {
	ldapService = NewLDAPService(&config.LDAP)
	jwtService = NewJWTService(&config.JWT)
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AuthProvider string `json:"authProvider"`
}

// loginHandler handles user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if loginReq.Username == "" || loginReq.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Check if user exists in database
	existingUser, err := GetUserByUsername(loginReq.Username)

	var user *User
	var authError error

	if err == nil && existingUser != nil {
		// User exists, determine auth method
		if existingUser.AuthProvider == "ldap" || existingUser.LDAPEnabled {
			// LDAP authentication
			user, authError = authenticateLDAP(loginReq.Username, loginReq.Password)
		} else {
			// Local authentication
			user, authError = authenticateLocal(existingUser, loginReq.Password)
		}
	} else {
		// User doesn't exist, try LDAP (auto-provision)
		user, authError = authenticateLDAP(loginReq.Username, loginReq.Password)
	}

	if authError != nil {
		log.Printf("Authentication failed for user %s: %v", loginReq.Username, authError)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := jwtService.GenerateToken(user)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set secure cookie (use Secure=true in production with HTTPS)
	jwtService.SetTokenCookie(w, token, false)

	// Return response
	response := LoginResponse{
		User: UserResponse{
			ID:           user.ID.String(),
			Username:     user.Username,
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			AuthProvider: user.AuthProvider,
		},
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// logoutHandler handles user logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear the auth token cookie
	jwtService.ClearTokenCookie(w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// getMeHandler returns the current authenticated user
func getMeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context (set by AuthMiddleware)
	user, err := GetUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := UserResponse{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AuthProvider: user.AuthProvider,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// authenticateLDAP authenticates against LDAP and syncs user
func authenticateLDAP(username, password string) (*User, error) {
	// Authenticate with LDAP
	ldapUser, err := ldapService.Authenticate(username, password)
	if err != nil {
		return nil, err
	}

	// Sync user to database
	user, err := SyncUserFromLDAP(ldapUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// authenticateLocal authenticates against local database
func authenticateLocal(user *User, password string) (*User, error) {
	// Verify password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
