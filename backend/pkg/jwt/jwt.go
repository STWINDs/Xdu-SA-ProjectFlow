package jwt

import (
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Claims is the custom JWT claims struct, embedding the standard registered claims.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwtv5.RegisteredClaims
}

// GenerateAccessToken creates a signed JWT access token with full user identity claims.
func GenerateAccessToken(userID uint, username, role string, secret string, ttl int) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(now.Add(time.Duration(ttl) * time.Second)),
			IssuedAt:  jwtv5.NewNumericDate(now),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken creates a signed JWT refresh token with only the user ID as subject.
func GenerateRefreshToken(userID uint, secret string, ttl int) (string, error) {
	now := time.Now()
	claims := jwtv5.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwtv5.NewNumericDate(now.Add(time.Duration(ttl) * time.Second)),
		IssuedAt:  jwtv5.NewNumericDate(now),
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken parses and validates a JWT token string, returning the custom Claims on success.
func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(t *jwtv5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwtv5.ErrSignatureInvalid
	}
	return claims, nil
}
