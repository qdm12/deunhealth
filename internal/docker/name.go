package docker

import "github.com/docker/docker/api/types"

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
