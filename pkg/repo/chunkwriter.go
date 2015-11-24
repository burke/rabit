package repo

/*
Copyright 2011 Google Inc.
Modifications Copyright 2015 Shopify Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/burke/rabit/pkg/rollsum"
)

const (
	// maxBlobSize is the largest blob we ever make when cutting up
	// a file.
	maxBlobSize = 1 << 20

	// bufioReaderSize is an explicit size for our bufio.Reader,
	// so we don't rely on NewReader's implicit size.
	// We care about the buffer size because it affects how far
	// in advance we can detect EOF from an io.Reader that doesn't
	// know its size.  Detecting an EOF bufioReaderSize bytes early
	// means we can plan for the final chunk.
	bufioReaderSize = 32 << 10

	// tooSmallThreshold is the threshold at which rolling checksum
	// boundaries are ignored if the current chunk being built is
	// smaller than this.
	tooSmallThreshold = 64 << 10
)

// noteEOFReader keeps track of when it's seen EOF, but otherwise
// delegates entirely to r.
type noteEOFReader struct {
	r      io.Reader
	sawEOF bool
}

func (r *noteEOFReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == io.EOF {
		r.sawEOF = true
	}
	return
}

type span struct {
	br string
}

func sha1FromString(s string) string {
	s1 := sha1.New()
	s1.Write([]byte(s))
	return hex.EncodeToString(s1.Sum(nil))
}

func uploadString(repo Repo, br, chunk string) error {
	pth := repo.ChunkPath(br)
	_ = os.Mkdir(path.Dir(pth), 0755)
	return ioutil.WriteFile(pth, []byte(chunk), 0660)
}

type chunkWriter struct {
	path  string
	r     io.Reader
	spans []span
}

func newChunkWriter(cspath string, r io.Reader) *chunkWriter {
	return &chunkWriter{path: cspath, r: r}
}

func (w *chunkWriter) writeChunks(repo Repo) ([]span, error) {
	var outerr error
	src := &noteEOFReader{r: w.r}
	bufr := bufio.NewReaderSize(src, bufioReaderSize)
	w.spans = []span{} // the tree of spans, cut on interesting rollsum boundaries
	rs := rollsum.New()
	var buf bytes.Buffer
	blobSize := 0 // of the next blob being built, should be same as buf.Len()

	const chunksInFlight = 32 // at ~64 KB chunks, this is ~2MB memory per file
	gate := make(chan struct{}, chunksInFlight)
	firsterrc := make(chan error, 1)

	// uploadLastSpan runs in the same goroutine as the loop below and is responsible for
	// starting uploading the contents of the buf.  It returns false if there's been
	// an error and the loop below should be stopped.
	uploadLastSpan := func() bool {
		chunk := buf.String()
		buf.Reset()
		select {
		case outerr = <-firsterrc:
			return false
		default:
			// No error seen so far, continue.
		}
		gate <- struct{}{}
		idx := len(w.spans) - 1
		go func() {
			defer func() { <-gate }()
			br := sha1FromString(chunk)
			w.spans[idx].br = br
			if err := uploadString(repo, br, chunk); err != nil {
				select {
				case firsterrc <- err:
				default:
				}
			}
		}()
		return true
	}

	for {
		c, err := bufr.ReadByte()
		if err != nil {
			if err == io.EOF {
				if blobSize > 0 {
					w.spans = append(w.spans, span{})
					if !uploadLastSpan() {
						return nil, outerr
					}
				}
				break
			} else {
				return nil, err
			}
		}

		buf.WriteByte(c)
		blobSize++
		onRollSplit := rs.Roll(c)
		switch {
		case blobSize == maxBlobSize || onRollSplit && blobSize > tooSmallThreshold:
			// split
		case src.sawEOF:
			// Don't split. End is coming soon enough.
			continue
		default:
			// Don't split.
			continue
		}
		blobSize = 0

		w.spans = append(w.spans, span{})

		if !uploadLastSpan() {
			return nil, outerr
		}
	}

	// Loop was already hit earlier.
	if outerr != nil {
		return nil, outerr
	}

	// Wait for all uploads to finish, one way or another, and then
	// see if any generated errors.
	// Once this loop is done, we own all the tokens in gate, so nobody
	// else can have one outstanding.
	for i := 0; i < chunksInFlight; i++ {
		gate <- struct{}{}
	}
	select {
	case err := <-firsterrc:
		return nil, err
	default:
	}

	return w.spans, nil
}
