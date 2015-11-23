package main

import (
	"fmt"
	"os"

	"github.com/burke/rabit/Godeps/_workspace/src/github.com/flynn/go-docopt"

	"github.com/burke/rabit/pkg/repo"
)

func init() {
	register("init", cmdInit, false, false, `
usage: %s init

Initialize a new rabit repository.

Environment Variables:
  RABIT_DIR  the pre-existing empty directory at which to create the repo
`)
}

func cmdInit(args *docopt.Args, rabitDir, rabitRemote string) error {
	if rabitDir == "" {
		return fmt.Errorf("RABIT_DIR must specify a path to an existing directory")
	}
	stat, err := os.Stat(rabitDir)
	if err != nil || !stat.IsDir() {
		return fmt.Errorf("RABIT_DIR must specify a path to an existing directory")
	}

	return repo.New(rabitDir).Init()
}
