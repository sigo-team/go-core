package lib

import (
	"github.com/golang-jwt/jwt"
)

func NewToken(claims jwt.StandardClaims, secret string) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(secret))
}

func ParseToken(s string, secret string) (*jwt.StandardClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(s, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken.Claims.(*jwt.StandardClaims), err
}
