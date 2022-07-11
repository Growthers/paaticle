package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	//fmt.Println(string(body), res)
	buildOutPutParser(string(body))

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

type buildOutput struct {
	Stream string `json:"stream"`
	Aux    auxID  `json:"aux"`
}

type auxID struct {
	ID string `json:"ID"`
}

type parsedBuildOutput struct {
	Log     string
	ImageID string
}

/*
DockerのImage Buildの出力はJSON形式で、
	{"stream": "出力"}
	{"aux": {"ID": "sha256"}}

*/
func buildOutPutParser(Out string) (res parsedBuildOutput, err error) {
	var readLine []string

	// ビルドの出力を1行ごとに分解する
	reader := bufio.NewReader(strings.NewReader(Out))
	for {
		l, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		readLine = append(readLine, l)
	}

	var parsedOutPutStrings []string
	for _, v := range readLine {
		tmp := buildOutput{}
		// 出力はJSONなのでパース
		err := json.Unmarshal([]byte(v), &tmp)
		if err != nil {
			return
		}

		parsedOutPutStrings = append(parsedOutPutStrings, tmp.Stream)
		if tmp.Aux.ID != "" {
			// ビルドされたImageのIDの出力の先頭に sha256: がついているのでカット
			res.ImageID = tmp.Aux.ID[7:]
		}
	}

	// 出力を連結
	res.Log = strings.Join(parsedOutPutStrings, "")

	return
}
