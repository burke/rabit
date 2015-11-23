package main

import (
	"os"

	"github.com/burke/rabit/chunkstore"
)

func main() {

	cs := chunkstore.New("/tmp/omgwtf", "whatever")

	f, err := os.Open("/tmp/vmlinuz")
	if err != nil {
		panic(err)
	}

	err = cs.Commit(f, "vmlinuz")

	if err != nil {
		panic(err)
	}
}
