package main

import (
	"fmt"
	"log"
)

// SyncUserFromLDAP creates or updates a user from LDAP data
func SyncUserFromLDAP(ldapUser *LDAPUser) (*User, error) {
	var user User

	// Check if user exists by username
	result := DB.Where("username = ?", ldapUser.Username).First(&user)

	if result.Error != nil {
		// User doesn't exist, create new
		user = User{
			Username:     ldapUser.Username,
			Email:        ldapUser.Email,
			FirstName:    ldapUser.FirstName,
			LastName:     ldapUser.LastName,
			PasswordHash: "", // Empty for LDAP users
			AuthProvider: "ldap",
			LDAPEnabled:  true,
		}

		if err := DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		log.Printf("Created new LDAP user: %s (%s)", user.Username, user.Email)
		return &user, nil
	}

	// User exists, update with LDAP data
	user.Email = ldapUser.Email
	user.FirstName = ldapUser.FirstName
	user.LastName = ldapUser.LastName
	user.AuthProvider = "ldap"
	user.LDAPEnabled = true
	user.PasswordHash = "" // Clear password hash for LDAP users

	if err := DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("Updated existing user from LDAP: %s (%s)", user.Username, user.Email)
	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID string) (*User, error) {
	var user User
	if err := DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
