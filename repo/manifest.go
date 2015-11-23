package repo

import (
	"strings"
)

type Manifest struct {
	chunks []string
}

func (m *Manifest) String() string {
	return strings.Join(m.chunks, "\n") + "\n"
}
