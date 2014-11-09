package got

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// global config

type Got struct {
	Handlers handlersConf
}

type Handlers struct {
	Packages    string
	GuiPackages string `yaml:"gui_packages"`
}

// profile config

type Profile struct {
	Description string
	Packages    []string
}

// package config

type Package struct {
	Name           string `yaml:"package"`
	Disabled       bool
	InstallActions []Action `yaml:"install"`
	UpdateActions  []Action `yaml:"update"`
}

type actionConf struct {
	Link    map[string]string
	Command string `yaml:"cmd"`
}

func ParseGot(path string) (*Got, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	gc := Got{}
	err = yaml.Unmarshal(c, &gc)
	return &gc, err
}

func ParseProfile(path string) (*Profile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	pc := Profile{}
	err = yaml.Unmarshal(c, &pc)
	return &pc, err
}

func ParsePackage(path string) (*Package, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	pc := Package{}
	err = yaml.Unmarshal(c, &pc)
	return &pc, err
}
