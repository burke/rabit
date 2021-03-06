package main

import (
	"github.com/burke/rabit/Godeps/_workspace/src/github.com/flynn/go-docopt"

	"github.com/burke/rabit/pkg/repo"
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
	repo := repo.New(rabitDir)
	name := args.String["<name>"]
	return repo.Rm(name)
}
