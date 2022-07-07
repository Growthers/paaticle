package git

import (
	"fmt"
	"os/exec"
)

func Clone(RepositoryURL string, AppID string) (res []byte, err error) {
	command := fmt.Sprintf("/usr/bin/git clone %s ./%s", RepositoryURL, AppID)
	res, err = exec.Command("/bin/sh", "-c", command).CombinedOutput()
	if err != nil {
		return
	}
	return
}
