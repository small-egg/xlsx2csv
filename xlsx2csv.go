package xlsx2csv

import (
	"bytes"
	"encoding/csv"
	"io"

	"github.com/tealeg/xlsx/v2"
)

// XLSXReader implements the io.Reader interface
// row by row converting an XLSX sheet to CSV
type XLSXReader struct {
	Align     bool // Align empty fields on the end to header length (see with_empty_fields)
	headerLen int

	data *xlsx.Sheet

	row int // Current row

	buff   *bytes.Buffer
	writer *csv.Writer
}

// NewReader creates instance of XLSXReader for specified sheet
func NewReader(data []byte, getSheet SheetGetter, comma rune) (*XLSXReader, error) {
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, err
	}

	sheet, err := getSheet(file)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(nil)

	csvWriter := csv.NewWriter(buff)
	csvWriter.Comma = comma

	reader := &XLSXReader{
		data:   sheet,
		buff:   buff,
		writer: csvWriter,
	}

	return reader, nil
}

// Read writes comma-separated byte representation
// of next row in XLSX sheet to b
func (r *XLSXReader) Read(p []byte) (n int, err error) {
	// Read to the end of current row
	if r.buff.Len() != 0 {
		return r.buff.Read(p)
	}

	if r.row >= r.data.MaxRow {
		return 0, io.EOF
	}

	row, err := r.nextRow()
	if err != nil {
		return 0, err
	}

	// If the first row was just read (header must be in first row)
	if r.row == 1 {
		r.headerLen = len(row)
	} else if r.Align && len(row) < r.headerLen {
		row = append(row, make([]string, r.headerLen-len(row))...)
	}

	err = r.writer.Write(row)
	if err != nil {
		return 0, err
	}

	r.writer.Flush()

	return r.buff.Read(p)
}

func (r *XLSXReader) nextRow() ([]string, error) {
	var row *xlsx.Row
	for row == nil {
		if r.row >= r.data.MaxRow {
			return nil, io.EOF
		}

		var err error
		row, err = r.data.Row(r.row)
		if err != nil {
			return nil, err
		}
		r.row++
	}

	res := make([]string, 0, r.data.MaxCol)
	for i :=0; i < r.data.MaxCol; i++ {
		c := row.GetCell(i)
		val, err := c.FormattedValue()
			if err != nil {
				res = append(res, err.Error())
			} else {
				res = append(res, val)
			}
	}

	return res, nil
}
