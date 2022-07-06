package git

import (
	"fmt"
	"os/exec"
)

func Clone(RepositoryURL string, AppID string) (LocalRepositoryPath string, err error) {
	command := fmt.Sprintf("/usr/bin/git clone %s ./%s", RepositoryURL, AppID)
	err = exec.Command("/bin/sh", "-c", command).Start()
	if err != nil {
		return "", err
	}
	LocalRepositoryPath = ""
	return
}
