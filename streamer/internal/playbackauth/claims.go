package playbackauth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID  int `json:"user_id"`
	MovieID int `json:"movie_id"`
	jwt.RegisteredClaims
}
