package got

import (
	"fmt"
	"path"

	"github.com/mikazuki/got/osx/brew"
)

func Parse(basePath string) (*Got, error) {
	mainConf, err := conf.ParseGot(path.Join(basePath, "got.yaml"))
	if err != nil {
		return nil, err
	}

	g := Got{}

	if mainConf.Handlers.Packages == "brew" {
		g.PackageManager = brew.New()
	} else {
		return nil, fmt.Errorf("package manager '%s' not implementedd", mainConf.Handlers.Packages)
	}

	return &g, nil
}
