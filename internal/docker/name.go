package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
)

func extractName(container types.Container) (name string) {
	name = container.ID
	for _, containerName := range container.Names {
		if containerName != "" {
			name = containerName
			break
		}
	}
	return name
}

func extractNameFromActor(actor events.Actor) (name string) {
	return actor.Attributes["name"]
}

func extractImageFromActor(actor events.Actor) (image string) {
	return actor.Attributes["image"]
}
