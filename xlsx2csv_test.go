package xlsx2csv

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertXLSXToCSV(t *testing.T) {
	type testData struct {
		file  string
		align bool

		result  [][]string
		errText string
	}

	data := []testData{
		{
			file: "testfiles/simple.xlsx",
			result: [][]string{
				{"Header1", "Header2"},
				{"Cell1", "Cell2"},
				{"Cell3", "Cell4"},
			},
		},
		{
			file: "testfiles/with_comma.xlsx",
			result: [][]string{
				{"Header1", "Header2", "Header3"},
				{"1", "x,", "2"},
				{"y", "3", "z`"},
			},
		},
		{
			file:  "testfiles/with_empty_cells.xlsx",
			align: true,
			result: [][]string{
				{"Header1", "Header2"},
				{"one", "two"},
				{"three", ""},
				{"", "four"},
			},
		},
		{
			file:    "testfiles/with_empty_cells.xlsx",
			align:   false,
			errText: "record on line 3: wrong number of fields",
		},
	}

	assert := assert.New(t)

	for _, testCase := range data {
		rawXLSX := readFile(testCase.file, assert)

		reader, err := NewReader(rawXLSX, WithName("sheet"), ',')
		assert.NoError(err)
		reader.Align = testCase.align

		rawCSV, err := ioutil.ReadAll(reader)
		assert.NoError(err)

		csvReader := csv.NewReader(bytes.NewReader(rawCSV))
		records, err := csvReader.ReadAll()
		if len(testCase.errText) == 0 {
			assert.NoError(err)
			assert.Equal(testCase.result, records)
		} else {
			assert.EqualError(err, testCase.errText)
		}
	}
}

func readFile(path string, assert *assert.Assertions) []byte {
	file, err := os.Open(path)
	assert.NoError(err)
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	assert.NoError(err)

	return data
}
