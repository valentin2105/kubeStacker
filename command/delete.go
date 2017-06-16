package command

import (
	"fmt"
	"os/exec"

	"github.com/codegangsta/cli"
	flag "github.com/ogier/pflag"
)

var (
	isDeleteAll bool
)

func init() {
	flag.BoolVarP(&isDeleteAll, "all", "a", false, "Delete completly the stack and volumes")
}

func deleteStack() {
	helmPath := CatchEnvHelm()
	helmDeleteCMD := fmt.Sprintf("%s delete --purge %s ", helmPath, stackName)
	exec.Command("sh", "-c", helmDeleteCMD).Output()
}

func deleteStackAll() {
	// Delete stack template files
	deployTmplPath := getConfigKey("deployTmplPath")
	thisDeployPath := fmt.Sprintf("%s/%s", deployTmplPath, stackName)
	deleteDeployPath := fmt.Sprintf("rm -r %s", thisDeployPath)
	Run(deleteDeployPath)
	// umount the Volume
	mountPlace := getConfigKey("mountPlace")
	stackMD5 := GetMD5Hash(stackName)
	umountStackVolume := fmt.Sprintf("umount %s/%s", mountPlace, stackMD5)
	Run(umountStackVolume)
	Run("sleep 1")
	// Delete Logical Volume
	volumeGroup := getConfigKey("volumeGroup")
	removeLV := fmt.Sprintf("lvremove -f /dev/%s/%s", volumeGroup, stackMD5)
	Run(removeLV)
	// Append fstab line
	cleanFstabCMD := fmt.Sprintf("sed -i 's,/dev/mapper/%s-%s	%s/%s               btrfs    defaults 0  1,,g' /etc/fstab", volumeGroup, stackMD5, mountPlace, stackMD5)
	Run(cleanFstabCMD)
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
