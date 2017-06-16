package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	flag "github.com/ogier/pflag"
)

var (
	isDeleteAll bool
)

func init() {
	flag.BoolVarP(&isDeleteAll, "all", "a", false, "Delete completly the stack and volumes")
	flag.StringVarP(&stackName, "name", "n", "", "Stack Name")
}

func deleteStack() {
	helmPath := CatchEnvHelm()
	helmDeleteCMD := fmt.Sprintf("%s delete --purge %s ", helmPath, stackName)
	Run(helmDeleteCMD)
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
	deleteStack()
	if isDeleteAll == true {
		deleteStackAll()
	}

}
