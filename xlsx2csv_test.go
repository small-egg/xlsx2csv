package xlsx2csv

import (
	"bytes"
	"encoding/csv"
	"io"
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
		{
			file:  "testfiles/with_unicode.xlsx",
			align: true,
			result: [][]string{
				{"Header1", "Header2", "Header3"},
				{"a", "a", "Bodrum'un Ortakent bölgesinde denize 150 metre mesafede kurulmuş her-şey-dahil bir tesis olan Medisun Hotel plaj üzerinde özel"},
				{"b", "b", "a"},
			},
		},
	}

	assert := assert.New(t)

	for _, testCase := range data {
		rawXLSX := readFile(testCase.file, assert)
		defer rawXLSX.Close()
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

func readFile(path string, assert *assert.Assertions) io.ReadCloser {
	file, err := os.Open(path)
	assert.NoError(err)

	return file
}
