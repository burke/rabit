package repo

import (
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
	path string
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

	return writeManifest(c.manifestPath(name), spans)
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
	mp := c.manifestPath(name)
	if err := os.Remove(mp); err != nil {
		return err
	}
	return c.GC(false)
}

func (c *repo) ChunkPath(hash string) string {
	prefix := hash[0:2]
	return filepath.Join(c.path, "chunks", prefix, hash)
}

func (c *repo) manifestPath(name string) string {
	return filepath.Join(c.path, "manifests", name)
}

func (c *repo) loadManifest(name string) (*manifest, error) {
	p := c.manifestPath(name)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	hashes := strings.Split(strings.TrimSpace(string(data)), "\n")
	return &manifest{chunks: hashes}, nil
}

func writeManifest(path string, spans []span) error {
	var hashes []string
	for _, span := range spans {
		hashes = append(hashes, span.br)
	}

	manifest := manifest{chunks: hashes}

	return ioutil.WriteFile(path, []byte(manifest.String()), 0660)
}
