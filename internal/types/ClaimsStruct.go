package types

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
