package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikazuki/got"
)

func main() {
	app := cli.NewApp()
	app.Name = "got"
	app.Version = "0.1.0"
	app.Usage = "manage your dotfiles with go"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:      "init",
			ShortName: "i",
			Usage:     "installs packages, creates symlinks and generates environment",
			Action:    cliInit,
		},
		{
			Name:      "update",
			ShortName: "u",
			Usage:     "performs update actions on your packages",
			Action:    func(c *cli.Context) { println("not implemented yet!") },
		},
	}

	app.Run(os.Args)
}

const BasePath = "/Users/patrick/.dot-files-new"

func cliInit(c *cli.Context) {
	gc, err := got.Parse(BasePath)
	fmt.Printf("%#v (err %v)\n", gc, err)
}
