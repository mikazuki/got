package tree

import (
	"errors"
	"path"
	"strings"
)

type Got struct {
	PackageManager    string
	GuiPackageManager string
	Profiles          []Profile
}

type Profile struct {
	Path        string
	Description string
	Enabled     bool
	Packages    []Package
}

type Package struct {
	Path             string
	InstallerPackage string
	Enabled          bool
	InstallActions   []Action
	UpdateActions    []Action
}

type Action interface {
	RunPackage(pkg Package) error
}

type LinkAction struct {
	Source, Target string
}

func (l *LinkAction) RunPackage(pkg Package) error {
	return errors.New("not implemented yet")
}

type CommandAction struct {
	Command string
}

func (c *CommandAction) RunPackage(pkg Package) error {
	return errors.New("not implemented yet")
}

func (g *Got) EnabledProfiles() []string {
	ps := []string{}
	for _, p := range g.Profiles {
		if p.Enabled {
			ps = append(ps, p.Name())
		}
	}
	return ps
}

func (p *Profile) Name() string {
	return path.Base(p.Path)
}

func (p *Package) Name() string {
	return strings.Split(p.InstallerPackage, " ")[0]
}

func (p *Package) IsExtended() bool {
	return path.Base(p.Path) == p.Name()
}
