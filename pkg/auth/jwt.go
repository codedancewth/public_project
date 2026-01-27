package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Configurable via environment variables
// JWT_SECRET: HMAC secret
// JWT_EXPIRE_MINUTES: access token expiry in minutes (default 15)
func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	return []byte(secret)
}

func getAccessExpiry() time.Duration {
	if v := os.Getenv("JWT_EXPIRE_MINUTES"); v != "" {
		if d, err := time.ParseDuration(v + "m"); err == nil {
			return d
		}
	}
	return 15 * time.Minute
}

// Claims represents JWT claims used in this project
// Minimal: sub (user id), exp, iat, jti, device_id optional

type Claims struct {
	UserID   int64  `json:"user_id"`
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken issues a signed JWT for the given user
func GenerateAccessToken(userID int64, deviceID string) (string, time.Time, error) {
	expireAt := time.Now().Add(getAccessExpiry())
	claims := &Claims{
		UserID:   userID,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(getSecret())
	return signed, expireAt, err
}

// ParseAndValidate parses a JWT string and returns claims if valid
func ParseAndValidate(tokenString string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
