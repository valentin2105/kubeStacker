package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/jmoiron/jsonq"
	flag "github.com/ogier/pflag"
)

var (
	stackName  string
	stackType  string
	configPath string
	volumeSize int
)

// Init add flags
func init() {
	flag.StringVarP(&stackName, "name", "n", "", "Stack Name")
	flag.StringVarP(&stackType, "type", "t", "", "Stack Type (Wordpress/Drupal...)")
	flag.IntVarP(&volumeSize, "size", "s", 0, "Stack Size (in GB)")
}

func catchEnvConfig() string {
	configPath := os.Getenv("KST_CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}
	return configPath
}

func getConfigKey(configKey string) string {
	ConfigPath := catchEnvConfig()
	b, err := ioutil.ReadFile(ConfigPath) // just pass the file name
	Check(err)
	str := string(b) // convert content to a 'string'
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(str))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	brutJson, err := jq.String("config", configKey)
	configKeyStr := string(brutJson)
	return configKeyStr
}

// Create a volume on the host
func createVolume(volumeName string, volumeSize int) {
	maxVolumeSize := getConfigKey("maxVolumeSize")
	maxVolumeSizeInt, _ := strconv.ParseInt(maxVolumeSize, 10, 0)
	if int64(volumeSize) > maxVolumeSizeInt {
		fmt.Printf("The maximum volume size is %sGB", maxVolumeSize)
		os.Exit(1)
	}
	volumeSizeStr := strconv.Itoa(volumeSize)
	volumeType := getConfigKey("volumeType")
	if volumeType == "lvm" {
		fmt.Printf("Let's Add a volume called %s with size of %sGB\n", volumeName, volumeSizeStr)

		volumeGroup := getConfigKey("volumeGroup")
		lvCreateCmd := fmt.Sprintf("lvcreate -L +%sG -n %s %s", volumeSizeStr, volumeName, volumeGroup)
		// lvcreate
		Run(lvCreateCmd)
		formatBtrfsCmd := fmt.Sprintf("mkfs.btrfs -f /dev/%s/%s", volumeGroup, volumeName)
		// mkfs.btrfs
		Run(formatBtrfsCmd)
		mountPlace := getConfigKey("mountPlace")
		volumeMountPlace := fmt.Sprintf("%s/%s", mountPlace, volumeName)
		if _, err := os.Stat(volumeMountPlace); os.IsNotExist(err) {
			os.Mkdir(volumeMountPlace, 0775)
		}
		// add to fstab and mount volume
		fstabCmd := fmt.Sprintf("/dev/mapper/%s-%s	%s/%s               btrfs    defaults 0  1\n", volumeGroup, volumeName, mountPlace, volumeName)
		AppendStringToFile("/etc/fstab", fstabCmd)
		//Run("mount -a")

	} else {
		fmt.Printf("This volumeType is not currently supported.")
		os.Exit(1)
	}
}

func helmInstall() {
	//helmPath := CatchEnvHelm()
}

// Main() for add command
func CmdAdd(c *cli.Context) {
	flag.Parse()
	titles := color.New(color.FgWhite, color.Bold)
	stackMD5 := GetMD5Hash(stackName)
	stackPath := CheckStackPath(stackType)
	stackPathExist := CheckStackPathExist(stackPath)
	if stackPathExist == false {
		panic("The Stack folder from config.json doesn't exist")
	}
	fmt.Printf("\n")
	fmt.Printf("\n")
	titles.Printf("Let's add %s (%s) -> %s... \n", stackName, stackMD5, stackType)
	fmt.Printf("\n")
	createVolume(stackMD5, volumeSize)

}
