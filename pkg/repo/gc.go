package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TODO(burke): this is total crap. rewrite.
func (c *repo) GC(verbose bool) error {
	allChunks := make(map[string]struct{})

	manifestDir := filepath.Join(c.path, "manifests")
	fis, err := ioutil.ReadDir(manifestDir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		p := c.manifestPath(fi.Name())
		data, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		hashes := strings.Split(strings.TrimSpace(string(data)), "\n")
		for _, h := range hashes {
			allChunks[h] = struct{}{}
		}
	}

	prefixFIs, err := ioutil.ReadDir(filepath.Join(c.path, "chunks"))
	if err != nil {
		return err
	}
	for _, fi := range prefixFIs {
		fis, err := ioutil.ReadDir(filepath.Join(c.path, "chunks", fi.Name()))
		if err != nil {
			return err
		}
		for _, cfi := range fis {
			if _, ok := allChunks[cfi.Name()]; !ok {
				p := filepath.Join(c.path, "chunks", fi.Name(), cfi.Name())
				_ = p
				if verbose {
					fmt.Println(cfi.Name())
				}
				os.Remove(p)
			}
		}
		fis, err = ioutil.ReadDir(filepath.Join(c.path, "chunks", fi.Name()))
		if err == nil && len(fis) == 0 {
			os.Remove(filepath.Join(c.path, "chunks", fi.Name()))
		}
	}

	return nil
}
