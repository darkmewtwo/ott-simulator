package capability

import (
	"simulator/internal/platform/capability/authentication"
	"simulator/internal/platform/httpclient"
)

type Manager struct {
	Authentication *authentication.Capability
}

func New(client *httpclient.Client) *Manager {
	return &Manager{
		Authentication: authentication.New(),
	}
}
