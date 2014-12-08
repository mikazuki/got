package got

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

type Context struct {
	GotDir         string
	ActiveProfiles []string
	Log            *logrus.Logger
}

func NewWdContext() (*Context, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	ctx := Context{GotDir: wd, ActiveProfiles: envProfiles(), Log: log}
	return &ctx, nil
}

func envProfiles() []string {
	gp := os.Getenv("GOT_PROFILES")
	return strings.Split(gp, ":")
}
