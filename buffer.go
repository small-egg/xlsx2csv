package xml2csv

import "errors"

type writer struct {
	buf []byte

	curr int
}

func newWriter(target []byte) *writer {
	w := new(writer)
	w.Reset(target)

	return w
}

func (w *writer) Write(p []byte) (n int, err error) {
	if len(w.buf)-w.curr > len(p) {
		return 0, errors.New("Buffer overflow")
	}

	n = copy(w.buf[w.curr:], p)
	w.curr += n

	return n, nil
}

func (w writer) Len() int {
	return len(w.buf)
}

func (w *writer) Reset(target []byte) {
	w.buf = target
	w.curr = 0
}
