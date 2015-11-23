package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("fetch", cmdPush, true, true, `
usage: %s push <name>

Copy a file from the local rabit repository to the rabit server.

Environment Variables:
  RABIT_DIR     Path on disk to the rabit repository
  RABIT_REMOTE  URL of remote rabit repository
`)
}

func cmdPush(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done push")

	return nil
}
