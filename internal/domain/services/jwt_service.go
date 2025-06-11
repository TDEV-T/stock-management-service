package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService defines the interface for JWT operations.
type JWTService interface {
	GenerateToken(userID uint, username string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type jwtServiceImpl struct {
	secretKey []byte
	issuer    string
	expiry    time.Duration
}

// NewJWTService creates a new JWTService.
// secretKey should be loaded from a secure configuration.
// issuer is a string identifying the token issuer.
// expiryHours defines the token validity period in hours.
func NewJWTService(secretKey string, issuer string, expiryHours int) (JWTService, error) {
	if secretKey == "" {
		return nil, errors.New("jwt secret key cannot be empty")
	}
	return &jwtServiceImpl{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		expiry:    time.Hour * time.Duration(expiryHours),
	}, nil
}

func (s *jwtServiceImpl) GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"iss":      s.issuer,
		"exp":      time.Now().Add(s.expiry).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", errors.New("could not sign token: " + err.Error())
	}
	return tokenString, nil
}

func (s *jwtServiceImpl) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
