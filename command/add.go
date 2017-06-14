package command

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func catchEnvConfig() string {
	configPath := os.Getenv("KST_CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}
	return configPath
}

func catchEnvHelm() string {
	helmPath := os.Getenv("HELM_PATH")
	if helmPath == "" {
		helmPath = "/usr/local/bin/helm"
	}
	return helmPath
}

// Get MD5 from a string (stackName)
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Exec shell command
func Run(command string) {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("command failed: %s\n", command)
		panic(err)
	}
}

// Check if stack already exist
func CheckStackExist(stackPath string) bool {
	checkPath := Exists(stackPath)
	if checkPath == true {
		return true
	} else {
		return false
	}
}

// Check if stackPath exist
func CheckStackPathExist(stackPath string) bool {
	checkPath := Exists(stackPath)
	if checkPath == true {
		return true
	} else {
		return false
	}
}

// Get and parse StackType from config.file
func CheckStackPath(stackType string) string {
	ConfigPath := catchEnvConfig()
	b, err := ioutil.ReadFile(ConfigPath) // just pass the file name
	if err != nil {
		panic(err)
	}
	str := string(b) // convert content to a 'string'
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(str))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	brutJson, err := jq.String(stackType)
	strStackPath := string(brutJson)
	return strStackPath
}

func getConfigKey(configKey string) string {
	ConfigPath := catchEnvConfig()
	b, err := ioutil.ReadFile(ConfigPath) // just pass the file name
	if err != nil {
		panic(err)
	}
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
func CreateVolume(volumeName string, volumeSize int) {
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
		fileHandle, _ := os.OpenFile("/root/fstab", os.O_APPEND, 0666)
		writer := bufio.NewWriter(fileHandle)
		defer fileHandle.Close()

		fmt.Fprintln(writer, "String I want to append")
		writer.Flush()

	} else {
		fmt.Printf("This volumeType is not currently supported.")
		os.Exit(1)
	}
}

func helmInstall() {
	//helmPath := catchEnvHelm()
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
	CreateVolume(stackMD5, volumeSize)

}
