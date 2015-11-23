package main

import (
	"github.com/burke/rabit/Godeps/_workspace/src/github.com/flynn/go-docopt"

	"github.com/burke/rabit/pkg/repo"
)

func init() {
	register("gc", cmdGC, true, false, `
usage: %s gc

Remove any blocks belonging only to removed manifests

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdGC(args *docopt.Args, rabitDir, rabitRemote string) error {
	repo := repo.New(rabitDir)
	return repo.GC(true)
}
