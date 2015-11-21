package main

import (
	"os"

	"github.com/burke/rabin-tuf/chunkstore"
)

func main() {

	cs := chunkstore.New("/tmp/omgwtf", "whatever")

	f, err := os.Open("test/vmlinuz")
	if err != nil {
		panic(err)
	}

	err = cs.Commit(f, "vmlinuz")

	if err != nil {
		panic(err)
	}
}
