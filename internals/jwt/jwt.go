package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService defines the interface for token operations.
type TokenService interface {
	GenerateToken(username string) (string, error)
	Validate(tokenString string) (bool, error)
	GenerateRefreshToken() (string, error)
	HashToken(token string) string
}

type JWT struct {
	Secret []byte
}

// NewJWT creates a new JWT instance with the provided secret.
func NewJWT(secret string) *JWT {
	return &JWT{Secret: []byte(secret)}
}

func (s *JWT) GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString(s.Secret)
}

func (s *JWT) Validate(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.Secret, nil
	})
	if err != nil {
		return false, err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}
	return false, nil
}

func (s *JWT) GenerateRefreshToken() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
	if err != nil {
		return "", err
	}
    return base64.URLEncoding.EncodeToString(b), nil
}

func (s *JWT) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}