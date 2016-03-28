package main

import (
	"os"

	"github.com/burke/rabit/pkg/manifest"
)

func main() {
	f, err := os.Open("/Users/burke/libv8-3.16.14.13.gem")
	if err != nil {
		panic(err)
	}
	m, err := manifest.Generate(f)
	if err != nil {
		panic(err)
	}
	m.Dump()
}
