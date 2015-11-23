package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("fetch", cmdFetch, `
usage: %s fetch <name>

Copy a file from a remote rabit server to the local repository.

Environment Variables:
  RABIT_DIR     Path on disk to the rabit repository
  RABIT_REMOTE  URL of remote rabit repository
`)
}

func cmdFetch(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done")

	return nil
}
