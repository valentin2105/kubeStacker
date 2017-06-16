package command

import (
	"fmt"
	"os/exec"

	"github.com/codegangsta/cli"
	flag "github.com/ogier/pflag"
)

var (
	isDeleteAll bool
	//stackName   string
)

func init() {
	//flag.StringVarP(&stackName, "name", "n", "", "Stack Name to Delete")
	flag.BoolVarP(&isDeleteAll, "all", "a", false, "Delete completly the stack and volumes")
}

func deleteStack() {
	helmPath := CatchEnvHelm()
	helmDeleteCMD := fmt.Sprintf("%s delete --purge %s ", helmPath, stackName)
	exec.Command("sh", "-c", helmDeleteCMD).Output()
}

func deleteStackAll() {
	deployTmplPath := getConfigKey("deployTmplPath")
	thisDeployPath := fmt.Sprintf("%s/%s", deployTmplPath, stackName)
	deleteDeployPath := fmt.Sprintf("rm -r %", thisDeployPath)
	Run(deleteDeployPath)
	mountPlace := getConfigKey("mountPlace")
	stackMD5 := GetMD5Hash(stackName)
	umountStackVolume := fmt.Sprintf("umount %s/%s", mountPlace, stackMD5)
	Run(umountStackVolume)
	volumeGroup := getConfigKey("volumeGroup")
	removeLV := fmt.Sprintf("lvremove -f /dev/%s/%s", volumeGroup, stackMD5)
	Run(removeLV)
}

func CmdDelete(c *cli.Context) {
	flag.Parse()
	if stackName == "" {
		panic("You need to give stack name (--name=)")
	}
	deleteStack()
	if isDeleteAll == true {
		deleteStackAll()
	}

}
