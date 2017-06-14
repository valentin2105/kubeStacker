package command

import (
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
	"github.com/jmoiron/jsonq"
	flag "github.com/ogier/pflag"
)

var (
	stackName  string
	stackType  string
	volumeSize int
)

// Init add flags
func init() {
	flag.StringVarP(&stackName, "name", "n", "", "Stack Name")
	flag.StringVarP(&stackType, "type", "t", "", "Stack Type (Wordpress/Drupal...)")
	flag.IntVarP(&volumeSize, "size", "s", 0, "Stack Size (in GB)")
}

// Get MD5 from a string (stackName)
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Exec shell command
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

// Check if stack already exist
func CheckStackExist() {

}

// Get and parse StackType from stacks.json
func CheckStackPath(stackType string) {
	b, err := ioutil.ReadFile("stacks.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(str))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	fmt.Println(jq.String(stackType))

}

// Create a LVM volume on the host
func CreateVolume(volumeName string, volumeSize int) {
	volumeSizeStr := strconv.Itoa(volumeSize)
	fmt.Printf("Let's Add a volume called %s with size of %sGB\n", volumeName, volumeSizeStr)
}

// Main() for add command
func CmdAdd(c *cli.Context) {
	flag.Parse()
	stackMD5 := GetMD5Hash(stackName)
	fmt.Printf("Let's add %s (%s) on %s... \n", stackName, stackMD5, stackType)
	fmt.Printf("\n")
	CreateVolume(stackMD5, volumeSize)
	CheckStackPath(stackType)

}
