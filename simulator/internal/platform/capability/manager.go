package capability

import (
	"simulator/internal/platform/capability/authentication"
	"simulator/internal/platform/capability/catalog"
	"simulator/internal/platform/capability/playback"
	"simulator/internal/platform/httpclient"
)

type Manager struct {
	Authentication *authentication.Capability
	Catalog        *catalog.Capability
	Playback       *playback.Capability
}

func New(client *httpclient.Client) *Manager {
	return &Manager{
		Authentication: authentication.New(),
		Catalog:        catalog.New(),
		Playback:       playback.New(),
	}
}
