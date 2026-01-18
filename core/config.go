package main

import (
	"log"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	LDAP LDAPConfig
	JWT  JWTConfig
}

// LDAPConfig holds LDAP-specific configuration
type LDAPConfig struct {
	URL          string
	BaseDN       string
	BindDN       string
	BindPassword string
	UserFilter   string
	UserBaseDN   string
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

var config *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	if config != nil {
		return config
	}

	expirationHours := 24 // default
	if hours := os.Getenv("JWT_EXPIRATION_HOURS"); hours != "" {
		if h, err := strconv.Atoi(hours); err == nil {
			expirationHours = h
		}
	}

	config = &Config{
		LDAP: LDAPConfig{
			URL:          getEnv("LDAP_URL", "ldap://openldap:389"),
			BaseDN:       getEnv("LDAP_BASE_DN", "dc=ironsnake,dc=local"),
			BindDN:       getEnv("LDAP_BIND_DN", "cn=admin,dc=ironsnake,dc=local"),
			BindPassword: getEnv("LDAP_BIND_PASSWORD", "admin"),
			UserFilter:   getEnv("LDAP_USER_FILTER", "(uid=%s)"),
			UserBaseDN:   getEnv("LDAP_USER_BASE_DN", "ou=users,dc=ironsnake,dc=local"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			ExpirationHours: expirationHours,
		},
	}

	// Validate required configuration
	if config.JWT.Secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return config
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	if config == nil {
		return LoadConfig()
	}
	return config
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
