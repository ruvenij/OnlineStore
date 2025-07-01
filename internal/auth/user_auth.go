package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserAuth struct {
	jwtSecret string
}

func NewUserAuth(secret string) *UserAuth {
	return &UserAuth{
		jwtSecret: secret,
	}
}

func (a *UserAuth) GenerateToken(userId string, userName string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userId,
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}
