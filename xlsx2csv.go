package xlsx2csv

import (
	"bytes"
	"encoding/csv"
	"io"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// XLSXReader implements the io.Reader interface
// row by row converting an XLSX sheet to CSV.
type XLSXReader struct {
	cfg config

	// Deprecated. Use WithAlign option instead.
	Align bool

	headerLen int

	data [][]string

	row int // Current row

	buff   *bytes.Buffer
	writer *csv.Writer
}

// New creates instance of XLSXReader
func New(reader io.Reader, options ...Option) (*XLSXReader, error) {
	file, err := excelize.OpenReader(reader)
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
	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	csvWriter := csv.NewWriter(buff)
	csvWriter.Comma = cfg.comma

	xlsxReader := &XLSXReader{
		data:   rows,
		buff:   buff,
		writer: csvWriter,
	}

	return xlsxReader, nil
}

// Deprecated. Use New instead
func NewReader(reader io.Reader, getSheet SheetSelector, comma rune) (*XLSXReader, error) {
	return New(reader, SetSheetSelector(getSheet), SetComma(comma))
}

// Read writes comma-separated byte representation
// of next row in XLSX sheet to b
func (r *XLSXReader) Read(p []byte) (n int, err error) {
	// Read to the end of current row
	if r.buff.Len() != 0 {
		return r.buff.Read(p)
	}

	if r.row >= len(r.data) {
		return 0, io.EOF
	}

	row, err := r.nextRow()
	if err != nil {
		return 0, err
	}

	// If the first row was just read (header must be in first row)
	if r.row == 1 {
		r.headerLen = len(row)
	} else if (r.cfg.align || r.Align) && len(row) < r.headerLen {
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
	var row []string
	for row == nil {
		if r.row >= len(r.data) {
			return nil, io.EOF
		}

		row = r.data[r.row]
		r.row++
	}

	res := make([]string, 0, len(row))
	for _, cell := range row {
		res = append(res, cell)

	}

	return res, nil
}
