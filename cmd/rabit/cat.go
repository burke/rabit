package main

import (
	"os"

	"github.com/burke/rabit/Godeps/_workspace/src/github.com/flynn/go-docopt"

	"github.com/burke/rabit/pkg/repo"
)

func init() {
	register("cat", cmdCat, true, false, `
usage: %s cat <name>

Write the contents of a file in the repository to stdout.

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdCat(args *docopt.Args, rabitDir, rabitRemote string) error {
	repo := repo.New(rabitDir, rabitRemote)

	name := args.String["<name>"]

	return repo.CatFile(name, os.Stdout)
}
