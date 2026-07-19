package session

import (
	"simulator/internal/platform/capability/authentication"
	"simulator/internal/user"
)

type Session struct {
	User                *user.User
	AuthenticationState authentication.AuthenticationState
}

func NewSession(user *user.User) *Session {
	return &Session{
		User: user,
	}
}
