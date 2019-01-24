package xml2csv

import "errors"

type writer struct {
	buf []byte

	curr int
	cap  int
}

func newWriter(target []byte) *writer {
	w := new(writer)
	w.Reset(target)

	return w
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.cap-w.curr >= len(p) {
		return 0, errors.New("Buffer overflow")
	}

	w.buf = append(w.buf, p...)
	w.curr += len(p)

	return len(p), nil
}

func (w writer) Len() int {
	return len(w.buf)
}

func (w *writer) Reset(target []byte) {
	w.cap = len(target)
	w.buf = target[:0]
	w.curr = 0
}
