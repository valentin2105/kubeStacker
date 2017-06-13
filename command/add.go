package command

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	flag "github.com/ogier/pflag"
)

var (
	stackName  string
	stackType  string
	volumeSize int
)

func init() {
	flag.StringVarP(&stackName, "name", "n", "", "Stack Name")
	flag.StringVarP(&stackType, "type", "t", "", "Stack Type (Wordpress/Drupal...)")
	flag.IntVarP(&volumeSize, "size", "s", 0, "Stack Size (in GB)")
}

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

func CreateVolume(volumeName string, volumeSize int) {
	volumeSizeStr := strconv.Itoa(volumeSize)
	fmt.Printf("Let's Add a volume called %s with size of %sGB\n", volumeName, volumeSizeStr)
}

func CmdAdd(c *cli.Context) {
	flag.Parse()
	volumeName := stackName
	fmt.Printf("Let's add %s with %s... \n", stackName, stackType)

	CreateVolume(volumeName, volumeSize)

}
