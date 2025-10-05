package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	IsOnboarded bool   `json:"isOnboarded"`
	jwt.RegisteredClaims
}

// TokenResponse represents the response structure for token operations
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// GenerateJWTToken generates a JWT token for the user
func GenerateJWTToken(userID, name, role string, isOnboarded bool) (string, error) {
	// Get JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET is not set")
    }
	claims := Claims{
		UserID:      userID,
		Name:        name,
		Role:        role,
		IsOnboarded: isOnboarded,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "espaze-backend",
			Subject:   userID,
		},
	}


	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWTToken validates a JWT token and returns the claims
func ValidateJWTToken(tokenString string) (*Claims, error) {
	// Get JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenExpirationTime returns the expiration time of a token
func GetTokenExpirationTime(tokenString string) (time.Time, error) {
	claims, err := ValidateJWTToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, nil
	}

	return claims.ExpiresAt.Time, nil
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(tokenString string) (bool, error) {
	expirationTime, err := GetTokenExpirationTime(tokenString)
	if err != nil {
		return true, err
	}

	if expirationTime.IsZero() {
		return false, nil
	}

	return time.Now().After(expirationTime), nil
}
