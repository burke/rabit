package chunkstore

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/burke/rabin-tuf/rollsum"
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
	from, to int64
	bits     int
	br       string
	children []span
}

func (s *span) isSingleBlob() bool {
	return len(s.children) == 0
}

func (s *span) size() int64 {
	size := s.to - s.from
	for _, cs := range s.children {
		size += cs.size()
	}
	return size
}

func sha1FromString(s string) string {
	s1 := sha1.New()
	s1.Write([]byte(s))
	return hex.EncodeToString(s1.Sum(nil))
}

func uploadString(br, chunk string) (n int, err error) {
	fmt.Printf("%s %d\n", string(br), len(chunk))
	return 0, nil
}

type chunkWriter struct {
	path  string
	r     io.Reader
	spans []span
}

func newChunkWriter(cspath string, r io.Reader) *chunkWriter {
	return &chunkWriter{path: cspath, r: r}
}

func (w *chunkWriter) writeChunks() (n int64, outerr error) {
	//func writeFileChunks(r io.Reader) (n int64, spans []span, outerr error) {
	src := &noteEOFReader{r: w.r}
	bufr := bufio.NewReaderSize(src, bufioReaderSize)
	w.spans = []span{} // the tree of spans, cut on interesting rollsum boundaries
	rs := rollsum.New()
	var last int64
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
		br := sha1FromString(chunk)
		w.spans[len(w.spans)-1].br = br
		select {
		case outerr = <-firsterrc:
			return false
		default:
			// No error seen so far, continue.
		}
		gate <- struct{}{}
		go func() {
			defer func() { <-gate }()
			if _, err := uploadString(br, chunk); err != nil {
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
		if err == io.EOF {
			if n != last {
				w.spans = append(w.spans, span{from: last, to: n})
				if !uploadLastSpan() {
					return
				}
			}
			break
		}
		if err != nil {
			return 0, err
		}

		buf.WriteByte(c)
		n++
		blobSize++
		rs.Roll(c)

		var bits int
		onRollSplit := rs.OnSplit()
		switch {
		case blobSize == maxBlobSize:
			bits = 20 // arbitrary node weight; 1<<20 == 1MB
		case src.sawEOF:
			// Don't split. End is coming soon enough.
			continue
		case onRollSplit && blobSize > tooSmallThreshold:
			bits = rs.Bits()
		default:
			// Don't split.
			continue
		}
		blobSize = 0

		// Take any spans from the end of the spans slice that
		// have a smaller 'bits' score and make them children
		// of this node.
		var children []span
		childrenFrom := len(w.spans)
		for childrenFrom > 0 && w.spans[childrenFrom-1].bits < bits {
			childrenFrom--
		}
		if nCopy := len(w.spans) - childrenFrom; nCopy > 0 {
			children = make([]span, nCopy)
			copy(children, w.spans[childrenFrom:])
			w.spans = w.spans[:childrenFrom]
		}

		w.spans = append(w.spans, span{from: last, to: n, bits: bits, children: children})
		last = n
		if !uploadLastSpan() {
			return
		}
	}

	// Loop was already hit earlier.
	if outerr != nil {
		return 0, outerr
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
		return 0, err
	default:
	}

	return n, nil
}
