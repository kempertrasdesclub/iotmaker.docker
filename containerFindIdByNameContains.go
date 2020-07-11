package iotmakerDocker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"strings"
)

func (el *DockerSystem) ContainerFindIdByNameContains(
	containsName string,
) (
	err error,
	list []NameAndId,
) {

	list = make([]NameAndId, 0)
	var listOfContainers []types.Container

	err, listOfContainers = el.ContainerListAll()
	for _, containerData := range listOfContainers {
		for _, containerName := range containerData.Names {
			if strings.Contains(containerName, containsName) == true {
				list = append(list, NameAndId{
					ID:   containerData.ID,
					Name: containerName,
				})
				return
			}
		}
	}

	err = errors.New("container name not found")

	return
}
