package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"stu-go/modules/user"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidTokenType   = errors.New("invalid token type")
)

type jwtService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	users      user.Service
}

// NewJWTService creates a JWT-backed auth service.
func NewJWTService(secret string, accessTTL, refreshTTL time.Duration, users user.Service) Service {
	return &jwtService{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		users:      users,
	}
}

func (s *jwtService) Login(creds Credentials) (TokenPair, error) {
	usr, err := s.users.GetByUsername(creds.Identifier)
	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	if usr.Password != creds.Password {
		return TokenPair{}, ErrInvalidCredentials
	}

	return s.generateTokens(usr.Username)
}

func (s *jwtService) Logout(_ string) error {
	// Stateless JWT logout is a no-op until token revocation is implemented.
	return nil
}

func (s *jwtService) Refresh(refreshToken string) (TokenPair, error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return TokenPair{}, err
	}
	if claims.Type != "refresh" {
		return TokenPair{}, ErrInvalidTokenType
	}
	return s.generateTokens(claims.Username)
}

func (s *jwtService) ValidateAccessToken(token string) (*Claims, error) {
	claims, err := s.parseToken(token)
	if err != nil {
		return nil, err
	}
	if claims.Type != "access" {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}

func (s *jwtService) generateTokens(username string) (TokenPair, error) {
	accessToken, err := s.createToken(username, "access", s.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := s.createToken(username, "refresh", s.refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *jwtService) createToken(username, tokenType string, ttl time.Duration) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	now := time.Now()
	payload := Claims{
		Username: username,
		Type:     tokenType,
		IssuedAt: now.Unix(),
		Expires:  now.Add(ttl).Unix(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)

	sig := sign(unsigned, s.secret)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(sig)

	return fmt.Sprintf("%s.%s", unsigned, signatureEncoded), nil
}

func (s *jwtService) parseToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	unsigned := fmt.Sprintf("%s.%s", parts[0], parts[1])
	expectedSig := sign(unsigned, s.secret)
	providedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid token signature")
	}

	if !hmac.Equal(expectedSig, providedSig) {
		return nil, errors.New("signature verification failed")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	claims := &Claims{}
	if err := json.Unmarshal(payloadBytes, claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	if time.Now().Unix() > claims.Expires {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func sign(data string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}
