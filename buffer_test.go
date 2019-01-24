package xml2csv

import (
	"math/rand"
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

func TestWriterOverflow(t *testing.T) {
	assert := assert.New(t)

	target := make([]byte, 1)
	writer := newWriter(target)

	data := []byte("01")
	n, err := writer.Write(data)
	assert.EqualError(err, "Buffer overflow")
	assert.Equal(0, n)
}

func TestWriterFewRecords(t *testing.T) {
	assert := assert.New(t)

	maxChunkSize := 256
	chunksCount := rand.Intn(10)

	target := make([]byte, maxChunkSize*chunksCount)
	writer := newWriter(target)

	res := make([]byte, 0, maxChunkSize*chunksCount)
	for i := 0; i < chunksCount; i++ {
		chunk := make([]byte, rand.Intn(maxChunkSize))
		size, _ := rand.Read(chunk)

		n, err := writer.Write(chunk)
		assert.NoError(err)
		assert.Equal(size, n)

		res = append(res, chunk...)
	}

	target = target[:len(res)]
	assert.Equal(res, target)
}
