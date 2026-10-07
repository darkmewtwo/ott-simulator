package session

import (
	"log"
	"simulator/internal/platform/capability/authentication"
	"simulator/internal/platform/httpclient"
	"simulator/internal/user"
	"time"
)

type Session struct {
	User           *user.User
	RegisterFirst  bool
	HTTPClient     *httpclient.Client
	authentication *authentication.Capability
}

func NewSession(user *user.User, registerFirst bool) *Session {
	client := httpclient.New("http://api:8100")
	return &Session{
		User:          user,
		RegisterFirst: registerFirst,
		HTTPClient:    client,
	}
}

func (s *Session) Run() error {
	log.Println("SESSION: ", s.RegisterFirst, s.User.Identity.Username)
	log.Println("USER STATE", s.User.MentalState)
	if s.RegisterFirst {
		if _, err := s.authentication.Register(s.HTTPClient, authentication.RegisterRequest{
			Username: s.User.Identity.Username,
			Password: s.User.Identity.Password,
			Email:    s.User.Identity.Email,
		}); err != nil {
			return err
		}
	}
	if _, err := s.authentication.Login(s.HTTPClient, authentication.LoginRequest{
		Username: s.User.Identity.Username,
		Password: s.User.Identity.Password,
	}); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	defer func() {
		if err := s.authentication.Logout(s.HTTPClient); err != nil {
			// TODO: log error
		}
	}()

	// TODO:
	// Start mental state routine
	// Start decision engine
	// Wait until session ends

	return nil
}
