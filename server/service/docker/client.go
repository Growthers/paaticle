package docker

import "github.com/docker/docker/client"

var cli *client.Client

func NewDockerClient() (err error) {
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	return nil
}
