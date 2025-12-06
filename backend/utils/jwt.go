package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// tạo secret key
func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	return []byte(secret)
}

// tạo token chứa id và role của user
func GenerateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(), // hết hạn sau 1 ngày
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

// xác minh token và trả về claims
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

// GenerateAccessToken tạo access token ngắn hạn (15 phút)
func GenerateAccessToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":         userID,
		"role":       role,
		"token_type": "access",
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

// GenerateRefreshToken tạo refresh token dài hạn (7 ngày)
func GenerateRefreshToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":         userID,
		"role":       role,
		"token_type": "refresh",
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

// ValidateRefreshToken kiểm tra refresh token hợp lệ
func ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}
	// Kiểm tra token type phải là "refresh"
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}
