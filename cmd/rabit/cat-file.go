package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("cat-file", cmdCatFile, `
usage: %s cat-file <name>

Write the contents of a file in the repository to stdout.

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdCatFile(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done")

	return nil
}
