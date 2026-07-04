package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret []byte
}

func New(jwtSecret []byte) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) ValidateToken(
	tokenString string,
) error {

	_, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {

			// Ensure the token uses HMAC signing
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return s.jwtSecret, nil
		},
	)

	if err != nil {
		return err
	}

	return nil
}
