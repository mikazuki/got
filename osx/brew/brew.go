package brew

import (
	"errors"

	"github.com/mikazuki/got"
)

type brewManager struct {
}

func New() got.PackageInstallerUpdater {
	return &brewManager{}
}

func (brew *brewManager) Install(pkg got.Package) error {
	return errors.New("not implemented")
}

func (brew *brewManager) Update(pkg got.Package) error {
	return errors.New("not implemented")
}
