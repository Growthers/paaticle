package docker

import (
	"context"
	"github.com/docker/docker/api/types"
)

func StartContainer(ContainerID string) error {
	err := cli.ContainerStart(context.Background(), ContainerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}
