package repo

import (
	"strings"
)

type manifest struct {
	chunks []string
}

func (m *manifest) String() string {
	return strings.Join(m.chunks, "\n") + "\n"
}
