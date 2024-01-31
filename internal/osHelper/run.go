package osHelper

import (
	"bytes"
	"fmt"
	"os/exec"
	"sb/internal/util"
)

func Run(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		util.LogErr(fmt.Sprintf("Error: %v \nstderr: %s", err, stderrBuf.String()))
	}
	util.LogInfoLn(stdoutBuf.String())
}
