package main

import (
	"github.com/flynn/go-docopt"

	"github.com/burke/rabit/chunkstore"
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
	repo := chunkstore.New(rabitDir, rabitRemote)

	name := args.String["<name>"]

	return repo.Rm(name)
}
