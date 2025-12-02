package auth

// Credentials captures user login input.
type Credentials struct {
	Identifier string
	Password   string
}

// TokenPair represents access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Service defines authentication operations.
type Service interface {
	Login(creds Credentials) (TokenPair, error)
	Logout(token string) error
	Refresh(refreshToken string) (TokenPair, error)
}
