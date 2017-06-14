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
	Result   string
)

func CatchEnvKubectl() string {
	kubectlPath := os.Getenv("KUBECTL_PATH")
	if kubectlPath == "" {
		kubectlPath = "/usr/local/bin/kubectl@"
	}
	return kubectlPath
}

// Exec shell command
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

// Check if a file exist
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Main() funct for Show
func CmdShow(c *cli.Context) {

	if len(os.Args) > 2 {
		var stackName = os.Args[2]
		stackMD5 := GetMD5Hash(stackName)
		var readyCmd = fmt.Sprintf("kubectl get deploy,pod,svc,secret,ingress -n %s", stackMD5)
		checkBin = Exists("/usr/local/bin/kubectl")

		if checkBin == true {
			fmt.Printf("Let's show %s (%s)\n", stackName, stackMD5)
			RunShow(readyCmd)

		} else {
			fmt.Printf("Kubectl is not present in /usr/local/bin\n")
		}

	} else {
		fmt.Println("Stack namespace is empty!")
		os.Exit(1)
	}
}
