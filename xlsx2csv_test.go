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
		exist bool

		result  [][]string
		errText string
	}

	data := []testData{
		{
			file:  "testfiles/simple.xlsx",
			exist: true,
			result: [][]string{
				{"Header1", "Header2"},
				{"Cell1", "Cell2"},
				{"Cell3", "Cell4"},
			},
		},
		{
			file:  "testfiles/with_comma.xlsx",
			exist: true,
			result: [][]string{
				{"Header1", "Header2", "Header3"},
				{"1", "x,", "2"},
				{"y", "3", "z`"},
			},
		},
		{
			file:  "testfiles/with_empty_cells.xlsx",
			align: true,
			exist: true,
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
			exist:   true,
			errText: "record on line 3: wrong number of fields",
		},
		{
			file:    "testfiles/sheet_not_exist.xlsx",
			align:   false,
			exist:   false,
			errText: "sheet sheet is not exist",
		},
		{
			file:  "testfiles/with_unicode.xlsx",
			align: true,
			exist: true,
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
		options := []Option{
			SetSheetSelector(SheetByName("sheet")),
			SetComma(','),
		}
		if testCase.align {
			options = append(options, WithAlign())
		}

		reader, err := New(rawXLSX, options...)
		if testCase.exist == false {
			assert.EqualError(err, testCase.errText, testCase.file)
			continue
		}
		assert.NoError(err)

		rawCSV, err := ioutil.ReadAll(reader)
		assert.NoError(err, testCase.file)

		csvReader := csv.NewReader(bytes.NewReader(rawCSV))
		records, err := csvReader.ReadAll()
		if len(testCase.errText) == 0 {
			assert.NoError(err, testCase.file)
			assert.Equal(testCase.result, records)
		} else {
			assert.EqualError(err, testCase.errText, testCase.file)
		}
	}
}

func readFile(path string, assert *assert.Assertions) io.ReadCloser {
	file, err := os.Open(path)
	assert.NoError(err)

	return file
}
