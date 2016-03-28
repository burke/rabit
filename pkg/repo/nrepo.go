package repo

import (
	"os"
	"path/filepath"
)

type NRepo interface {
	GenerateManifest(string) error
}

type nrepo struct {
	path string
}

func (n *nrepo) GenerateManifest(path string) error {
	fp := filepath.Join(n.path, path)
	r, err := os.Open(fp)
	if err != nil {
		return err
	}

	w := newChunkWriter(c.path, r)
	return nil
}
