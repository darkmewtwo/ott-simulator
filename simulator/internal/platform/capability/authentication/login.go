package authentication

import "simulator/internal/platform/httpclient"

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (c *Capability) Login(
	client *httpclient.Client,
	request LoginRequest,
) (*LoginResponse, error) {

	var response LoginResponse

	err := client.Post("/auth/login", request, &response)
	if err != nil {
		return nil, err
	}
	client.AccessToken = response.AccessToken
	client.TokenType = response.TokenType

	return &response, nil
}

func (c *Capability) Register(
	client *httpclient.Client,
	request RegisterRequest,
) (*RegisterResponse, error) {

	var response RegisterResponse

	err := client.Post("/auth/register", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Capability) Logout(
	client *httpclient.Client,
) error {
	client.AccessToken = ""
	client.TokenType = ""
	return nil
}
