package authentication

import httpclient "simulator/internal/platform/httpclient"

type Authentication struct {
	client *httpclient.Client
	state  State
}

func New() *Authentication {
	return &Authentication{}
}
