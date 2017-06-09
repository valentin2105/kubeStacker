package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
  "github.com/codegangsta/cli"
	)


func RunAdd(command string) {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("command failed: %s\n", command)
		panic(err)
	}
}


func CmdAdd(c *cli.Context) {
	RunAdd("echo Lets add a stack...")
}
