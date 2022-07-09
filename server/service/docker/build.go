package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io/ioutil"
	"log"
	"os"
)

func BuildDockerImage(ImageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	buildOpts := types.ImageBuildOptions{
		Dockerfile: ImageName,
		Tags:       []string{"paaticle"},
	}

	res, err := cli.ImageBuild(
		context.Background(),
		makeArchivedDockerfile(ImageName),
		buildOpts,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	return nil
}

// DockerはtarでまとめられたDockerfileしか受け取らないのでtarでまとめる
func makeArchivedDockerfile(fileName string) *bytes.Reader {
	f, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Panic(err)
		}
	}()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	// archive the Dockerfile
	tarHeader := &tar.Header{
		Name: fileName,
		Size: int64(len(b)),
	}
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Panic(err)
	}
	_, err = tw.Write(b)
	if err != nil {
		log.Panic(err)
	}

	return bytes.NewReader(buf.Bytes())
}
