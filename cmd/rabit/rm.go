package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("rm", cmdRm, true, false, `
usage: %s rm <name>

Remove a file from the rabit repository.
You may want to run 'rabit gc' afterward.

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdRm(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done rm")

	return nil
}
