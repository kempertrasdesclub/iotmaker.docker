package iotmakerDocker

import (
	"github.com/docker/docker/client"
)

// Negotiate best docker version
func (el *DockerSystem) ClientCreate() error {
	var err error

	el.cli, err = client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
		client.WithHTTPHeaders(
			map[string]string{
				"Content-Type": "application/tar",
			},
		),
	)

	return err
}
