package xlsx2csv

import (
	"bytes"
	"encoding/csv"
	"io"

	"github.com/tealeg/xlsx/v2"
)

// XLSXReader implements the io.Reader interface
// row by row converting an XLSX sheet to CSV.
type XLSXReader struct {
	cfg config

	// Deprecated. Use WithAlign option instead.
	Align bool

	headerLen int

	data *xlsx.Sheet

	row int // Current row

	buff   *bytes.Buffer
	writer *csv.Writer
}

// New creates instance of XLSXReader
func New(raw []byte, options ...Option) (*XLSXReader, error) {
	file, err := xlsx.OpenBinary(raw)
	if err != nil {
		return nil, err
	}

	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	sheet, err := cfg.getSheet(file)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(nil)

	csvWriter := csv.NewWriter(buff)
	csvWriter.Comma = cfg.comma

	reader := &XLSXReader{
		data:   sheet,
		buff:   buff,
		writer: csvWriter,
	}

	return reader, nil
}

// Deprecated. Use New instead
func NewReader(data []byte, getSheet SheetSelector, comma rune) (*XLSXReader, error) {
	return New(data, SetSheetSelector(getSheet), SetComma(comma))
}

// Read writes comma-separated byte representation
// of next row in XLSX sheet to b
func (r *XLSXReader) Read(p []byte) (n int, err error) {
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

	switch {
	case r.row == 1: // If the first row was just read (header must be in first row)
		r.headerLen = len(row)
	case (r.cfg.align || r.Align) && len(row) < r.headerLen:
		row = append(row, make([]string, r.headerLen-len(row))...)
	case len(row) > r.headerLen:
		row = row[:r.headerLen]
	}

	if err := r.writer.Write(row); err != nil {
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

		row = r.data.Row(r.row)
		r.row++
	}

	res := make([]string, 0, len(row.Cells))
	for _, cell := range row.Cells {
		val, err := cell.FormattedValue()
		if err != nil {
			res = append(res, err.Error())
		} else {
			res = append(res, val)
		}
	}

	return res, nil
}
