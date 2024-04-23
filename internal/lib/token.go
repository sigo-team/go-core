package lib

import (
	"github.com/golang-jwt/jwt"
)

const secret = "supersecretsecret"

func NewToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(secret))
}

func ParseToken(token string) (*jwt.StandardClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return parsedToken.Claims.(*jwt.StandardClaims), err
}
