package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
  "github.com/codegangsta/cli"
	flag "github.com/ogier/pflag"
)

var (
	stackName string
	stackType string
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
	flag.Parse()

	fmt.Printf("Let's add %s with %s... \n", stackName, stackType)


}


func init() {
  flag.StringVarP(&stackName, "name", "n", "", "Stack Name")
	flag.StringVarP(&stackType, "type", "t", "", "Stack Type (Wordpress/Drupal...)")
}
