package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/alifsuryadi/ecolokal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateJWT(userID int, email, role string, cfg *config.Config) (string, error) {
    expireHours, _ := strconv.Atoi(cfg.JWTExpire)
    expirationTime := time.Now().Add(time.Duration(expireHours) * time.Hour)

    claims := &JWTClaim{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateJWT(tokenString string, cfg *config.Config) (*JWTClaim, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(cfg.JWTSecret), nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*JWTClaim)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}