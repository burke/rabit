package main

import (
	"os"

	"github.com/flynn/go-docopt"

	"github.com/burke/rabit/chunkstore"
)

func init() {
	register("cat-file", cmdCatFile, true, false, `
usage: %s cat-file <name>

Write the contents of a file in the repository to stdout.

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdCatFile(args *docopt.Args, rabitDir, rabitRemote string) error {
	repo := chunkstore.New(rabitDir, rabitRemote)

	name := args.String["<name>"]

	return repo.CatFile(name, os.Stdout)
}
