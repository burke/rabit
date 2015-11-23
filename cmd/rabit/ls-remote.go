package main

import (
	"fmt"

	"github.com/burke/rabit/Godeps/_workspace/src/github.com/flynn/go-docopt"
)

func init() {
	register("ls-remote", cmdLsRemote, false, true, `
usage: %s ls-remote

List files in a remote rabit repository.

Environment Variables:
  RABIT_REMOTE  URL of remote rabit repository
`)
}

func cmdLsRemote(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done ls-remote")

	return nil
}
