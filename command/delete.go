package command

import (
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

}

func deleteStackAll() {

}

func CmdDelete(c *cli.Context) {
	flag.Parse()
	deleteStack()
	if isDeleteAll == true {
		deleteStackAll()
	}

}
