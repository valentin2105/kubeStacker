package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
  "github.com/codegangsta/cli"
	)

var (
	checkBin bool
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

func Exists(name string) bool {
    _, err := os.Stat(name)
    return !os.IsNotExist(err)
}

func CmdShow(c *cli.Context) {

	if len(os.Args) > 2 {
	    var ns = os.Args[2]
	    var readyCmd = fmt.Sprintf("kubectl get deploy,pod,svc,secret,ingress -n %s", ns)
			checkBin = Exists("/usr/local/bin/kubectl")
			if checkBin == true{

			  RunShow(readyCmd)

			} else{
			fmt.Printf("Kubectl is not present in /usr/local/bin\n")
			}

	} else {
			fmt.Println("Stack namespace is empty!")
			os.Exit(1)
	}
}
