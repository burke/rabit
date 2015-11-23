package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("init", cmdInit, `
usage: %s init

Initialize a new rabit repository.

Environment Variables:
  RABIT_DIR  the pre-existing empty directory at which to create the repo
`)
}

func cmdInit(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done")

	return nil
}
