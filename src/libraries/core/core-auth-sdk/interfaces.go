package core_auth_sdk

import (
	jose "gopkg.in/square/go-jose.v2"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

// Provides a JSON Web Key from a Key ID
// Wanted to use function signature from go-jose.v2
// but that would make us lose error information
type JWKProvider interface {
	Key(kid string) ([]jose.JSONWebKey, error)
}

// Extracts verified in-built claims from a jwt idToken
type JWTClaimsExtractor interface {
	GetVerifiedClaims(idToken string) (*jwt.Claims, error)
}

// AuthService exposes the interface contract the authentication service client adheres to
type AuthService interface {
	// Get a user account
	GetAccount(id string) (*Account, error)
	// Updates the username associated with a user account
	Update(id, username string) error
	// Locks a user account
	LockAccount(id string) error
	// Unlocks a user account
	UnlockAccount(id string) error
	// Archives a user account
	ArchiveAccount(id string) error
	// Creates a new user account
	ImportAccount(username, password string, locked bool) (int, error)
	// Expires the password associated with a user account
	ExpirePassword(id string) error
	// Authenticates a user account
	LoginAccount(username, password string) (string, error)
	// Signs up a user account
	SignupAccount(username, password string) (string, error)
	// Remove a session associated with a given user account
	LogOutAccount() error
}
