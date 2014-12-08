package tree

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/mikazuki/got"
	"github.com/mikazuki/got/conf"
)

const gotDir = "got"
const mainConfName = "got.yaml"
const profileConfName = "profile.yaml"
const pkgConfName = "package.yaml"

var ErrNoRootConf = errors.New("Root configuration not found")
var ErrNoProfileConf = errors.New("No profile configuration found")
var ErrNoPackageConf = errors.New("No package configuration found")

var log = logrus.New()

func Parse(ctx *got.Context) (*Got, error) {
	mainConf, err := conf.ParseGot(path.Join(ctx.GotDir, mainConfName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoRootConf
		}
		return nil, err
	}

	g := Got{
		PackageManager:    mainConf.Handlers.Packages,
		GuiPackageManager: mainConf.Handlers.GuiPackages,
		Profiles:          make([]Profile, 0, 8),
	}

	mainDir, err := os.Open(ctx.GotDir)
	if err != nil {
		return nil, err
	}

	subDirs, err := mainDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, subDir := range subDirs {
		logDir := ctx.Log.WithFields(logrus.Fields{
			"dir": subDir.Name(),
		})
		if !subDir.IsDir() {
			logDir.Debug("file encountered")
			continue
		}
		if subDir.Name()[0] == '.' {
			logDir.Debug("skipping hidden directory")
			continue
		}
		if subDir.Name() == gotDir {
			logDir.Debug("skipping got directory")
			continue
		}

		profile, err := parseProfile(path.Join(ctx.GotDir, subDir.Name()), ctx)
		if err != nil {
			logDir.WithFields(logrus.Fields{
				"err": err,
			}).Warn("could not parse profile directory")
			continue
		}

		g.Profiles = append(g.Profiles, *profile)
	}

	return &g, nil
}

func parseProfile(profilePath string, ctx *got.Context) (*Profile, error) {
	profileConf, err := conf.ParseProfile(path.Join(profilePath, profileConfName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoProfileConf
		}
		return nil, err
	}

	p := Profile{
		Path:        profilePath,
		Description: profileConf.Description,
		Enabled:     false,
		Packages:    make([]Package, len(profileConf.Packages)),
	}

	for _, pn := range ctx.ActiveProfiles {
		if pn == p.Name() {
			p.Enabled = true
		}
	}

	for i, pkgDef := range profileConf.Packages {
		p.Packages[i] = Package{
			Path:             profilePath,
			InstallerPackage: pkgDef,
			Enabled:          true,
		}
	}

	profileDir, err := os.Open(profilePath)
	if err != nil {
		return nil, err
	}

	subDirs, err := profileDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, subDir := range subDirs {
		logDir := ctx.Log.WithFields(logrus.Fields{
			"profile": profilePath,
			"dir":     subDir.Name(),
		})
		if !subDir.IsDir() {
			logDir.Debug("file encountered")
			continue
		}
		if subDir.Name()[0] == '.' {
			logDir.Debug("skipping hidden directory")
			continue
		}

		pkg, err := parsePackage(path.Join(profilePath, subDir.Name()), ctx)
		if err != nil {
			logDir.WithFields(logrus.Fields{
				"err": err,
			}).Warn("could not parse package directory")
			continue
		}

		p.Packages = append(p.Packages, *pkg)
	}

	return &p, nil
}

func parsePackage(pkgPath string, ctx *got.Context) (*Package, error) {
	pkgConf, err := conf.ParsePackage(path.Join(pkgPath, pkgConfName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoPackageConf
		}
		return nil, err
	}

	pkg := Package{
		Path:             pkgPath,
		InstallerPackage: pkgConf.Name,
		Enabled:          !pkgConf.Disabled,
		InstallActions:   make([]Action, 0, len(pkgConf.InstallActions)),
		UpdateActions:    make([]Action, 0, len(pkgConf.UpdateActions)),
	}

	for _, instDef := range pkgConf.InstallActions {
		actions, err := parseAction(&instDef)
		if err != nil {
			ctx.Log.WithFields(logrus.Fields{
				"packagePath": pkgPath,
				"err":         err,
			}).Warn("failed to parse action")
		}

		for _, action := range actions {
			pkg.InstallActions = append(pkg.InstallActions, action)
		}
	}

	return &pkg, nil
}

func parseAction(action *conf.Action) ([]Action, error) {
	if action.Command != "" && (action.Link != nil || len(action.Link) != 0) {
		return nil, fmt.Errorf("Encountered an action that is both link and command")
	}

	if action.Command != "" {
		return []Action{&CommandAction{Command: action.Command}}, nil
	} else if action.Link != nil {
		links := make([]Action, len(action.Link))
		i := 0
		for s, t := range action.Link {
			links[i] = &LinkAction{Source: s, Target: t}
			i++
		}
		return links, nil
	} else {
		return nil, fmt.Errorf("Unknown action encountered")
	}
}
