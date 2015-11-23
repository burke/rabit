package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("add", cmdInit, `
usage: %s add <path> <name>

Add a file to the rabit repository

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdAdd(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done")

	return nil
}
