package playbackauth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
}

func New(secret string) *Service {
	return &Service{
		secret: []byte(secret),
	}
}

func (s *Service) Validate(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return s.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
