package command

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jmoiron/jsonq"
)

// Check error
func Check(e error) {
	if e != nil {
		panic(e)
	}
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

// Write to a file func
func AppendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
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
	Check(err)
	str := string(b) // convert content to a 'string'
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(str))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	brutJson, err := jq.String(stackType)
	strStackPath := string(brutJson)
	return strStackPath
}

// Catch Helm Path from env
func CatchEnvHelm() string {
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
