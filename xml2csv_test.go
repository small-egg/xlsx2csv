package xml2csv

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertXLSXToCSV(t *testing.T) {
	data := map[string][][]string{
		"testfiles/simple.xlsx": {
			{"Header1", "Header2"},
			{"Cell1", "Cell2"},
			{"Cell3", "Cell4"},
		},
		"testfiles/with_comma.xlsx": {
			{"Header1", "Header2", "Header3"},
			{"1", "x,", "2"},
			{"y", "3", "z`"},
		},
	}

	assert := assert.New(t)

	for file, result := range data {
		rawXLSX := readFile(file, assert)

		reader, err := NewReader(rawXLSX, "sheet", ',')
		assert.NoError(err)

		rawCSV, err := ioutil.ReadAll(reader)
		assert.NoError(err)

		csvReader := csv.NewReader(bytes.NewReader(rawCSV))
		records, err := csvReader.ReadAll()
		assert.NoError(err)

		assert.Equal(result, records)
	}
}

func readFile(path string, assert *assert.Assertions) []byte {
	file, err := os.Open(path)
	assert.NoError(err)

	data, err := ioutil.ReadAll(file)
	assert.NoError(err)

	return data
}
