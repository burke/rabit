package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("gc", cmdGC, true, false, `
usage: %s gc

Remove any blocks belonging only to removed manifests

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdGC(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done gc")

	return nil
}
