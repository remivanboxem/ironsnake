package main

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// LDAPService handles LDAP authentication operations
type LDAPService struct {
	config *LDAPConfig
}

// LDAPUser represents a user retrieved from LDAP
type LDAPUser struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	DN        string
}

// NewLDAPService creates a new LDAP service instance
func NewLDAPService(config *LDAPConfig) *LDAPService {
	return &LDAPService{
		config: config,
	}
}

// Connect establishes a connection to the LDAP server
func (s *LDAPService) Connect() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(s.config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	return conn, nil
}

// Authenticate validates user credentials against LDAP and returns user info
func (s *LDAPService) Authenticate(username, password string) (*LDAPUser, error) {
	// Sanitize username to prevent LDAP injection
	username = sanitizeLDAPInput(username)

	conn, err := s.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Bind with admin credentials
	err = conn.Bind(s.config.BindDN, s.config.BindPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to bind with admin credentials: %w", err)
	}

	// Search for user
	searchFilter := fmt.Sprintf(s.config.UserFilter, username)
	searchRequest := ldap.NewSearchRequest(
		s.config.UserBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"uid", "mail", "givenName", "sn", "cn"},
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search for user: %w", err)
	}

	if len(searchResult.Entries) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	if len(searchResult.Entries) > 1 {
		return nil, fmt.Errorf("multiple users found with same username")
	}

	userEntry := searchResult.Entries[0]
	userDN := userEntry.DN

	// Attempt to bind as the user to verify password
	err = conn.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Extract user attributes
	ldapUser := &LDAPUser{
		DN:        userDN,
		Username:  username,
		Email:     userEntry.GetAttributeValue("mail"),
		FirstName: userEntry.GetAttributeValue("givenName"),
		LastName:  userEntry.GetAttributeValue("sn"),
	}

	// Fallback to cn if givenName/sn not available
	if ldapUser.FirstName == "" && ldapUser.LastName == "" {
		cn := userEntry.GetAttributeValue("cn")
		parts := strings.SplitN(cn, " ", 2)
		ldapUser.FirstName = parts[0]
		if len(parts) > 1 {
			ldapUser.LastName = parts[1]
		}
	}

	return ldapUser, nil
}

// GetUserAttributes retrieves user attributes from LDAP without authentication
func (s *LDAPService) GetUserAttributes(username string) (*LDAPUser, error) {
	// Sanitize username to prevent LDAP injection
	username = sanitizeLDAPInput(username)

	conn, err := s.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Bind with admin credentials
	err = conn.Bind(s.config.BindDN, s.config.BindPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to bind with admin credentials: %w", err)
	}

	// Search for user
	searchFilter := fmt.Sprintf(s.config.UserFilter, username)
	searchRequest := ldap.NewSearchRequest(
		s.config.UserBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"uid", "mail", "givenName", "sn", "cn"},
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search for user: %w", err)
	}

	if len(searchResult.Entries) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	userEntry := searchResult.Entries[0]

	// Extract user attributes
	ldapUser := &LDAPUser{
		DN:        userEntry.DN,
		Username:  username,
		Email:     userEntry.GetAttributeValue("mail"),
		FirstName: userEntry.GetAttributeValue("givenName"),
		LastName:  userEntry.GetAttributeValue("sn"),
	}

	// Fallback to cn if givenName/sn not available
	if ldapUser.FirstName == "" && ldapUser.LastName == "" {
		cn := userEntry.GetAttributeValue("cn")
		parts := strings.SplitN(cn, " ", 2)
		ldapUser.FirstName = parts[0]
		if len(parts) > 1 {
			ldapUser.LastName = parts[1]
		}
	}

	return ldapUser, nil
}

// sanitizeLDAPInput sanitizes user input to prevent LDAP injection
func sanitizeLDAPInput(input string) string {
	// Remove potentially dangerous characters
	replacer := strings.NewReplacer(
		"\\", "\\5c",
		"*", "\\2a",
		"(", "\\28",
		")", "\\29",
		"\x00", "\\00",
	)
	return replacer.Replace(input)
}
