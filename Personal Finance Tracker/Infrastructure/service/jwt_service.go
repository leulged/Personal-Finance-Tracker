package services

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "fmt"
)

type JWTService struct {
    SecretKey string
}

func NewJWTService(secretKey string) *JWTService {
    return &JWTService{
        SecretKey: secretKey,
    }
}

type JWTServiceInterface interface {
    GenerateTokens(userID string, userRole string) (string, string, error)
    ValidateToken(tokenString string) (jwt.MapClaims, error) 
    GenerateResetToken(userID string) (string, error)
    ValidateResetToken(tokenString string) (jwt.MapClaims, error)
}

// GenerateTokens creates access and refresh tokens with essential claims.
func (s *JWTService) GenerateTokens(userID string, userRole string) (string, string, error) {
    now := time.Now()

    // Access Token (15 min expiry)
    accessClaims := jwt.MapClaims{
        "sub": userID,
        "role": userRole,
        "iat": now.Unix(),
        "exp": now.Add(15 * time.Minute).Unix(),
        "jti": uuid.NewString(),
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString([]byte(s.SecretKey))
    if err != nil {
        return "", "", err
    }

    // Refresh Token (7 day expiry)
    refreshClaims := jwt.MapClaims{
        "sub": userID,
        "role": userRole,
        "iat": now.Unix(),                          // timestamp 
        "exp": now.Add(7 * 24 * time.Hour).Unix(),
        "jti": uuid.NewString(),
    }
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(s.SecretKey))
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
}

// ValidateToken verifies the token signature and expiration.
func (s *JWTService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(s.SecretKey), nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

// GenerateResetToken creates a reset token with 15-minute expiry
func (s *JWTService) GenerateResetToken(userID string) (string, error) {
    now := time.Now()

    resetClaims := jwt.MapClaims{
        "sub": userID,
        "type": "reset",
        "iat": now.Unix(),
        "exp": now.Add(15 * time.Minute).Unix(),
        "jti": uuid.NewString(),
    }
    
    resetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, resetClaims)
    resetTokenString, err := resetToken.SignedString([]byte(s.SecretKey))
    if err != nil {
        return "", err
    }

    return resetTokenString, nil
}

// ValidateResetToken verifies the reset token signature and expiration
func (s *JWTService) ValidateResetToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(s.SecretKey), nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid reset token")
    }

    // Check if it's a reset token
    if tokenType, ok := claims["type"].(string); !ok || tokenType != "reset" {
        return nil, fmt.Errorf("invalid token type")
    }

    return claims, nil
}
