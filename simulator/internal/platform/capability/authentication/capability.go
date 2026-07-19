package authentication

import httpclient "simulator/internal/platform/httpclient"

type Authentication struct {
	client *httpclient.Client
}

func New() *Authentication {
	return &Authentication{}
}
