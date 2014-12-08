package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/mikazuki/got"
	"github.com/mikazuki/got/tree"
)

var log = logrus.New()

func main() {
	app := cli.NewApp()
	app.Name = "got"
	app.Version = "0.1.0"
	app.Usage = "manage your dotfiles with go"
	app.Author = "Patrick Marschik"
	app.Email = "patrick@marschik.me"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "verbose output",
			EnvVar: "GOT_VERBOSE",
		},
		cli.BoolFlag{
			Name:   "quiet, q",
			Usage:  "quiet output",
			EnvVar: "GOT_QUIET",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initializes an empty got tree",
			Action: func(c *cli.Context) { println("not implemented yet!") },
		},
		{
			Name:      "install",
			ShortName: "i",
			Usage:     "installs packages, creates symlinks and generates environment",
			Action:    func(c *cli.Context) { println("not implemented yet!") },
		},
		{
			Name:      "update",
			ShortName: "u",
			Usage:     "performs update actions on your packages",
			Action:    func(c *cli.Context) { println("not implemented yet!") },
		},
		{
			Name:      "tree",
			ShortName: "t",
			Usage:     "shows complete active got tree",
			Action:    printTree,
		},
	}
	app.Run(os.Args)
}

func printTree(c *cli.Context) {
	ctx := makeContext(c)
	tree, err := tree.Parse(ctx)
	if err != nil {
		ctx.Log.WithFields(logrus.Fields{
			"err": err,
		}).Fatalf("Could not generate tree")
		return
	}

	fmt.Print("Enabled profiles: ")
	ep := tree.EnabledProfiles()
	for i, p := range ep {
		fmt.Print(p)
		if i != len(ep)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()

	fmt.Printf("package manager    : %s\n", tree.PackageManager)
	fmt.Printf("GUI package manager: %s\n", tree.GuiPackageManager)
	fmt.Println("active profiles:")
	hasActiveProfles := false
	for _, p := range tree.Profiles {
		if !p.Enabled {
			continue
		}

		hasActiveProfles = true

		fmt.Printf("  name       : %s\n", p.Name())
		fmt.Printf("  path       : %s\n", p.Path)
		fmt.Printf("  description: %s\n", p.Description)
		fmt.Println("  active packages:")

		hasActivePackages := false
		for _, pkg := range p.Packages {
			if !pkg.Enabled {
				continue
			}

			hasActivePackages = true

			fmt.Print("    ", pkg.Name())
			hasInstallFlags := pkg.Name() != pkg.InstallerPackage
			if hasInstallFlags || pkg.IsExtended() {
				fmt.Print(" (")

				if hasInstallFlags {
					fmt.Print(pkg.InstallerPackage)
					if pkg.IsExtended() {
						fmt.Print(", ")
					}
				}
				if pkg.IsExtended() {
					fmt.Printf("at %s", pkg.Path)
				}

				fmt.Printf(")")
			}
			fmt.Println()
		}
		if !hasActivePackages {
			fmt.Println("    none")
		}
	}
	if !hasActiveProfles {
		fmt.Println("  none")
	}
}

func makeContext(c *cli.Context) *got.Context {
	verbose := c.GlobalBool("verbose")
	quiet := c.GlobalBool("quiet")

	if verbose && quiet {
		log.Fatal("Flags quiet and verbose cannot be on at the same time")
		return nil
	}

	ctx, err := got.NewWdContext()
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("could not create working directory context")
		return nil
	}

	if verbose {
		ctx.Log.Level = logrus.DebugLevel
	} else if quiet {
		ctx.Log.Level = logrus.FatalLevel
	}

	return ctx
}
