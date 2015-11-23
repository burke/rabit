package repo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Repo interface {
	Init() error
	Add(io.Reader, string) error
	LsFiles() ([]string, error)
	CatFile(string, io.Writer) error
	Rm(string) error
	GC(bool) error
	ChunkPath(string) string
}

type repo struct {
	path   string
	remote string
}

func New(path string) Repo {
	return &repo{path: path}
}

func (c *repo) Init() error {
	if err := os.Mkdir(filepath.Join(c.path, "chunks"), 0755); err != nil {
		return err
	}

	return os.Mkdir(filepath.Join(c.path, "manifests"), 0755)
}

func (c *repo) Add(r io.Reader, name string) error {
	w := newChunkWriter(c.path, r)
	spans, err := w.writeChunks(c)
	if err != nil {
		return err
	}

	return writeManifest(c.ManifestPath(name), spans)
}

func (c *repo) CatFile(name string, w io.Writer) error {
	manifest, err := c.loadManifest(name)
	if err != nil {
		return err
	}

	for _, ch := range manifest.chunks {
		cpath := c.ChunkPath(ch)
		f, err := os.Open(cpath)
		if err != nil {
			return err
		}
		io.Copy(w, f)
		f.Close()
	}

	return nil
}

func (c *repo) LsFiles() ([]string, error) {
	manifestDir := filepath.Join(c.path, "manifests")
	fis, err := ioutil.ReadDir(manifestDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, fi := range fis {
		names = append(names, fi.Name())
	}
	return names, nil
}

func (c *repo) Rm(name string) error {
	mp := c.ManifestPath(name)
	if err := os.Remove(mp); err != nil {
		return err
	}
	return c.GC(false)
}

func (c *repo) GC(verbose bool) error {
	allChunks := make(map[string]struct{})

	manifestDir := filepath.Join(c.path, "manifests")
	fis, err := ioutil.ReadDir(manifestDir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		p := c.ManifestPath(fi.Name())
		data, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		hashes := strings.Split(strings.TrimSpace(string(data)), "\n")
		for _, h := range hashes {
			//fmt.Println(h)
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

func (c *repo) ChunkPath(hash string) string {
	prefix := hash[0:2]
	return filepath.Join(c.path, "chunks", prefix, hash)
}

func (c *repo) ManifestPath(name string) string {
	return filepath.Join(c.path, "manifests", name)
}

func (c *repo) loadManifest(name string) (*Manifest, error) {
	p := c.ManifestPath(name)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	hashes := strings.Split(strings.TrimSpace(string(data)), "\n")
	return &Manifest{chunks: hashes}, nil
}

func writeManifest(path string, spans []span) error {
	var hashes []string
	for _, span := range spans {
		hashes = append(hashes, span.br)
	}

	manifest := Manifest{chunks: hashes}

	return ioutil.WriteFile(path, []byte(manifest.String()), 0660)
}
