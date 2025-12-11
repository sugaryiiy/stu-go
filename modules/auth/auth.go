package auth

// Credentials captures user login input.
type Credentials struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// TokenPair represents access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Claims encodes JWT payload details.
type Claims struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	IssuedAt int64  `json:"iat"`
	Expires  int64  `json:"exp"`
}

// Service defines authentication operations.
type Service interface {
	Login(creds Credentials) (TokenPair, error)
	Logout(token string) error
	Refresh(refreshToken string) (TokenPair, error)
	ValidateAccessToken(token string) (*Claims, error)
}
