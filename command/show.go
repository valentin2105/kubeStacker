package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
  "github.com/codegangsta/cli"
	)

func RunShow(command string) {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("command failed: %s\n", command)
		panic(err)
	}
}

func CmdShow(c *cli.Context) {

	if len(os.Args) > 2 {
	    var arg = os.Args[2]
	    var readyCmd = fmt.Sprintf("kubectl get deploy,pod,svc,secret,ingress -n %s", arg)
			RunShow(readyCmd)

	} else {
			fmt.Println("Stack namespace is empty!")
			os.Exit(1)
	}


}
