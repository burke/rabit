package main

import (
	"fmt"

	"github.com/flynn/go-docopt"

	"github.com/burke/rabit/chunkstore"
)

func init() {
	register("ls-files", cmdLsFiles, true, false, `
usage: %s ls-files

List files in a rabit repository

Environment Variables:
  RABIT_DIR  Path on disk to the rabit repository
`)
}

func cmdLsFiles(args *docopt.Args, rabitDir, rabitRemote string) error {
	repo := chunkstore.New(rabitDir, rabitRemote)
	names, err := repo.LsFiles()
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}
