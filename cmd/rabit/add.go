package main

import (
	"os"

	"github.com/burke/rabit/Godeps/_workspace/src/github.com/flynn/go-docopt"

	"github.com/burke/rabit/pkg/repo"
)

func init() {
	register("add", cmdAdd, true, false, `
usage: %s add <path> <name>

Add a file to the rabit repository

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdAdd(args *docopt.Args, rabitDir, rabitRemote string) error {
	repo := repo.New(rabitDir)

	path := args.String["<path>"]
	name := args.String["<name>"]

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	return repo.Add(f, name)
}
