package got

import "errors"

type Got struct {
	PackageManager    PackageInstallerUpdater
	GuiPackageManager PackageInstallerUpdater
	Profiles          []Profile
}

type PackageInstaller interface {
	Install(pkg Package) error
}

type PackageUpdater interface {
	Update(pkg Package) error
}

type PackageInstallerUpdater interface {
	PackageInstaller
	PackageUpdater
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
	InstallActions   []InstallAction
	UpdateActions    []UpdateAction
}

type InstallAction interface {
	Install(pkg Package) error
}

type UpdateAction interface {
	Update(pkg Package) error
}

type LinkInstallAction struct {
	Source, Target string
}

func (l *LinkInstallAction) Install(pkg Package) error {
	return errors.New("not implemented yet")
}

type CommandUpdateAction struct {
	Command string
}

func (c *CommandUpdateAction) Update(pkg Package) error {
	return errors.New("not implemented yet")
}
