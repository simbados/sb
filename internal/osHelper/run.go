package osHelper

import (
	"os"
	"os/exec"
	"sb/internal/log"
)

func Run(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		log.LogErr(err)
	}
	if err := cmd.Wait(); err != nil {
		log.LogErr(err)
	}
}
