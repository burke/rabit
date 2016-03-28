package manifest

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

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

func sha1FromString(s string) string {
	s1 := sha1.New()
	s1.Write([]byte(s))
	return hex.EncodeToString(s1.Sum(nil))
}

type span struct {
	sum string
	off uint64
	len uint64
}

type Manifest struct {
	sum   string
	spans []span
}

func (m *Manifest) Dump() {
	for _, sp := range m.spans {
		fmt.Printf("%v\t%d\t%s\n", sp.off, sp.len, sp.sum)
	}
	fmt.Printf("\x1b[34m%s\x1b[0m\n", m.sum)
}

func Generate(r io.Reader) (*Manifest, error) {
	m := &Manifest{}
	fileSum := sha256.New()

	src := &noteEOFReader{r: r}
	bufr := bufio.NewReaderSize(src, bufioReaderSize)
	m.spans = []span{} // the tree of spans, cut on interesting rollsum boundaries
	rs := rollsum.New()
	var buf bytes.Buffer
	blobSize := 0 // of the next blob being built, should be same as buf.Len()

	var offset uint64

	appendSpan := func() {
		chunk := buf.String()
		buf.Reset()
		m.spans = append(m.spans, span{
			sum: sha1FromString(chunk),
			off: offset,
			len: uint64(len(chunk)),
		})
		offset += uint64(len(chunk))
	}

	for {
		c, err := bufr.ReadByte()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			// EOF
			if blobSize > 0 {
				appendSpan()
			}
			break
		}

		fileSum.Write([]byte{c})
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
		appendSpan()
	}

	m.sum = hex.EncodeToString(fileSum.Sum(nil))
	return m, nil
}
