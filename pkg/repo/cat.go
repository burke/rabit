package repo

import (
	"bufio"
	"io"
	"io/ioutil"
)

func (c *repo) CatFile(name string, w io.Writer) error {
	manifest, err := c.loadManifest(name)
	if err != nil {
		return err
	}

	// read up to 32 files in parallel.
	// This brings cat time for a 5GB file from 32s to 9s on my machine.
	files := make(chan *fileContentsOptionPromise, 32)

	go func() {
		defer close(files)

		for _, ch := range manifest.chunks {
			cpath := c.ChunkPath(ch)

			// when we have either data or an error, we resolve the promise
			of := fileContentsOptionPromise{resolved: make(chan struct{})}
			go func(of *fileContentsOptionPromise) {
				f, err := ioutil.ReadFile(cpath)
				if err != nil {
					of.err = err
				} else {
					of.data = f
				}
				close(of.resolved)
			}(&of)

			files <- &of
		}
	}()

	wb := bufio.NewWriter(w)

	for file := range files {
		<-file.resolved // wait until the promise has resolved to a value.
		if file.err != nil {
			return err
		}
		wb.Write(file.data)
	}

	wb.Flush()

	return nil
}
