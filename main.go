package main

import (
	"os"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "kubeStacker"
	app.Version = "v0.1"
  app.Author = "valentin2105"
	app.Email = "valentin@ouvrard.it"
	app.Usage = "Kubernetes Stack Deployer"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)

}
