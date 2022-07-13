package docker

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"os"
	"os/exec"
)

func CreateContainer(ImageID string, ContainerName string, Env []string, AppID string) (containerID string, warnings []string, err error) {
	containerConfig := &container.Config{
		Env:   Env,
		Image: ImageID,
	}

	create, err := cli.ContainerCreate(context.Background(), containerConfig, &container.HostConfig{}, nil, nil, ContainerName)
	if err != nil {
		return
	}

	_, err = makeTarBall(AppID)
	if err != nil {
		return
	}

	err = sendCodeToContainer(cli, create.ID, AppID)
	if err != nil {
		return
	}

	return create.ID, create.Warnings, nil
}

func sendCodeToContainer(cli *client.Client, ContainerID string, AppID string) error {
	archive, _ := os.Open(AppID + ".tar")
	defer archive.Close()

	return cli.CopyToContainer(context.Background(), ContainerID, "/", bufio.NewReader(archive), types.CopyToContainerOptions{})
}

/*
	コンテナにファイルをコピーするにはtarアーカイブで無ければいけないのでtarにまとめる
	.git/は除外
*/
func makeTarBall(AppID string) (Log string, err error) {
	command := fmt.Sprintf("cd %s;tar --exclude .git -cf ../%s.tar *", AppID, AppID)
	l, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		return
	}

	Log = string(l)
	return
}
