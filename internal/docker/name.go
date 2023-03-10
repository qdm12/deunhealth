package docker

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
)

func extractName(container types.Container) (name string) {
	name = container.ID
	for _, containerName := range container.Names {
		if containerName != "" {
			name = strings.TrimPrefix(containerName, "/")
			break
		}
	}
	return name
}

func extractNameFromActor(actor events.Actor) (name string) {
	name = actor.ID
	if actor.Attributes["name"] != "" {
		name = actor.Attributes["name"]
	}

	return name
}
