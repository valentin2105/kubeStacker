package command

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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
		Run("mount -a")
	} else {
		fmt.Printf("This volumeType is not currently supported.")
		os.Exit(1)
	}
}

func copyHelmTemplate(stackPath string) {
	deployTmplPath := getConfigKey("deployTmplPath")
	thisDeployPath := fmt.Sprintf("%s/%s", deployTmplPath, stackName)
	err := Copy_folder(stackPath, thisDeployPath)
	if err != nil {
		log.Fatal(err)
	} else {
		titles := color.New(color.FgWhite, color.Bold)
		fmt.Printf("\n")
		titles.Printf("Helm configuration copied to %s \n", thisDeployPath)
	}

}

func parseHelmTemplate(from string, to string) {
	t, err := template.ParseFiles(from)
	Check(err)
	f, err := os.Create(to)
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	// Helm template config
	stackMD5 := GetMD5Hash(stackName)
	mountPlace := getConfigKey("mountPlace")
	volumeMountPlace := fmt.Sprintf("%s/%s", mountPlace, stackMD5)
	volumeMountPlaceDB := fmt.Sprintf("%s/db", volumeMountPlace)
	volumeMountPlaceWeb := fmt.Sprintf("%s/web", volumeMountPlace)
	stackPasswd := "aBigStr0ngPassw0rd"
	config := map[string]string{
		"siteURL":       stackName,
		"siteMD5":       stackMD5,
		"rootPassword":  stackPasswd,
		"volumePathDB":  volumeMountPlaceDB,
		"volumePathWeb": volumeMountPlaceWeb,
	}
	err = t.Execute(f, config)
	Check(err)
	f.Close()
}

func helmInstall(stackPath string) {
	helmPath := CatchEnvHelm()
	deployTmplPath := getConfigKey("deployTmplPath")
	thisDeployPath := fmt.Sprintf("%s/%s", deployTmplPath, stackName)
	helmInitCMD := fmt.Sprintf("%s install --name %s %s", helmPath, stackName, thisDeployPath)
	//fmt.Printf(helmInitCMD)
	Run(helmInitCMD)
}

func createNamespace(stackMD5 string) {
	createNsCmd := fmt.Sprintf("/usr/bin/kubectl create ns %s", stackMD5)
	exec.Command("sh", "-c", createNsCmd).Output()
	//Run(createNsCmd)
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
	//Start Output
	fmt.Printf("\n")
	fmt.Printf("\n")
	titles.Printf("Let's add a %s for %s (%s) \n", stackType, stackName, stackMD5)
	fmt.Printf("\n")
	// Call volume creation
	createVolume(stackMD5, volumeSize)
	// Copy & Parse Helm Template
	deployTmplPath := getConfigKey("deployTmplPath")
	thisDeployPath := fmt.Sprintf("%s/%s", deployTmplPath, stackName)
	helmValueTmplPath := fmt.Sprintf("%s/values.tmpl.yaml", thisDeployPath)
	helmValuePath := fmt.Sprintf("%s/values.yaml", thisDeployPath)
	copyHelmTemplate(stackPath)
	parseHelmTemplate(helmValueTmplPath, helmValuePath)
	// Create k8s namespace
	createNamespace(stackMD5)
	// Install Helm generated package
	helmInstall(stackPath)
	fmt.Printf("\n")
	titles.Printf("https://%s is correctly deployed !\n", stackName)
}
