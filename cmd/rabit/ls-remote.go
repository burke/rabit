package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("ls-remote", cmdLsRemote, `
usage: %s ls-remote

List files in a remote rabit repository.

Environment Variables:
  RABIT_REMOTE  URL of remote rabit repository
`)
}

func cmdLsRemote(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done")

	return nil
}
