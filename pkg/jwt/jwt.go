package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// TokenClaims adalah struktur untuk claims JWT
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken menghasilkan token JWT baru
func GenerateToken(userID, email, role, secret string, expiryHours int) (string, error) {
	// Mengatur waktu kadaluarsa token
	expirationTime := time.Now().Add(time.Hour * time.Duration(expiryHours))

	// Membuat claims
	claims := &TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Membuat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan kunci rahasia
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi dan mem-parsing token JWT
func ValidateToken(tokenString, secret string) (*TokenClaims, error) {
	// Parse token dengan claims
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Mengambil claims dari token
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")
}

// RefreshToken memperbaharui token JWT lama dengan yang baru
func RefreshToken(tokenString, secret string, newExpiryHours int) (string, error) {
	// Validasi token lama
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return "", err
	}

	// Menghasilkan token baru dengan user_id dan email yang sama
	return GenerateToken(claims.UserID, claims.Email, claims.Role, secret, newExpiryHours)
}

// GetTokenRemainingValidity mengembalikan sisa waktu valid token dalam detik
func GetTokenRemainingValidity(tokenString, secret string) (int64, error) {
	// Validasi token
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return 0, err
	}

	// Hitung sisa waktu valid token
	now := time.Now().Unix()
	remaining := claims.ExpiresAt - now

	if remaining < 0 {
		return 0, errors.New("Token has expired")
	}

	return remaining, nil
}
