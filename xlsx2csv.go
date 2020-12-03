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
		cfg:    cfg,
		data:   rows,
		buff:   buff,
		writer: csvWriter,
	}

	return xlsxReader, nil
}

// Read writes comma-separated byte representation
// of next row in XLSX sheet to b
func (r *XLSXReader) Read(p []byte) (n int, err error) {
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
