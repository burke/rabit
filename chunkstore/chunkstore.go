package chunkstore

import (
	"io"
)

type ChunkStore interface {
	Commit(io.Reader, string) error
	Ls() ([]string, error)
	LsRemote() ([]string, error)
	Push(string) error
	Rm(string) error
	RmRemote(string) error
	Fetch(string) error
	GC() error
}

type chunkStore struct {
	path   string
	remote string
}

func New(path, remote string) ChunkStore {
	return &chunkStore{path: path, remote: remote}
}

func (c *chunkStore) Commit(r io.Reader, name string) error {
	w := newChunkWriter(c.path, r)
	_, err := w.writeChunks()
	return err
}

func (c *chunkStore) Ls() ([]string, error) {
	return nil, nil
}

func (c *chunkStore) LsRemote() ([]string, error) {
	return nil, nil
}

func (c *chunkStore) Push(name string) error {
	return nil
}

func (c *chunkStore) Rm(name string) error {
	return nil
}

func (c *chunkStore) RmRemote(name string) error {
	return nil
}

func (c *chunkStore) Fetch(name string) error {
	return nil
}

func (c *chunkStore) GC() error {
	return nil
}
