package git

import (
	"fmt"
	"os/exec"
)

func CheckOut(AppID string, BranchName string) (res []byte, err error) {
	command := fmt.Sprintf("cd ./%s ; /usr/bin/git switch %s ; /usr/bin/git branch", AppID, BranchName)
	res, err = exec.Command("/bin/sh", "-c", command).CombinedOutput()
	if err != nil {
		return
	}
	return
}
