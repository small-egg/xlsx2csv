package xml2csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	assert := assert.New(t)

	target := make([]byte, 10)
	writer := newWriter(target)

	data := []byte("0123456789")
	n, err := writer.Write(data)
	assert.NoError(err)
	assert.Equal(len(data), n)

	assert.Equal(data, target)
}
