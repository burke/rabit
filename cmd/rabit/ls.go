package main

import (
	"fmt"

	"github.com/burke/rabit/Godeps/_workspace/src/github.com/flynn/go-docopt"

	"github.com/burke/rabit/pkg/repo"
)

func init() {
	register("ls", cmdLs, true, false, `
usage: %s ls

List files in a rabit repository

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdLs(args *docopt.Args, rabitDir, rabitRemote string) error {
	repo := repo.New(rabitDir)
	names, err := repo.LsFiles()
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}
