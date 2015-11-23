package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
)

func init() {
	register("ls-files", cmdLsFiles, `
usage: %s ls-files

List files in a rabit repository

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdLsFiles(args *docopt.Args, rabitDir, rabitRemote string) error {
	fmt.Println("done")

	return nil
}
