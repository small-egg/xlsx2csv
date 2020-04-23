package xlsx2csv

import (
	"errors"
	"io"
	"io/ioutil"
)

type readerAt struct {
	r io.Reader
	offset int64
}

func newReaderAt(r io.Reader) io.ReaderAt {
	return &readerAt{r: r}
}

func (r *readerAt) ReadAt(p []byte, off int64) (int, error) {
	if off < r.offset {
		return 0, errors.New("invalid offset")
	}

	diff := off - r.offset
	written, err := io.CopyN(ioutil.Discard, r.r, diff)
	r.offset += written
	if err != nil {
		return 0, err
	}

	n, err := r.r.Read(p)
	r.offset += int64(n)
	return n, err
}
