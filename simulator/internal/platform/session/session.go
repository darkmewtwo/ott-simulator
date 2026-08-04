package session

import (
	"simulator/internal/platform/httpclient"
	"simulator/internal/user"
)

type Session struct {
	User       *user.User
	HTTPClient *httpclient.Client
}

func NewSession(user *user.User) *Session {
	client := httpclient.New("http://localhost:8100")
	return &Session{
		User:       user,
		HTTPClient: client,
	}
}
