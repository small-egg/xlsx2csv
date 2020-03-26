package xlsx2csv

import (
	"errors"

	"github.com/tealeg/xlsx/v2"
)

var (
	sheetNotFoundErr = errors.New("requested sheet not found")
)

type SheetGetter func(file *xlsx.File) (*xlsx.Sheet, error)

func WithName(name string) SheetGetter {
	return func(file *xlsx.File) (*xlsx.Sheet, error) {
		sheet, ok := file.Sheet[name]
		if !ok {
			return nil, sheetNotFoundErr
		}

		return sheet, nil
	}
}

func WithIndex(i int) SheetGetter {
	return func(file *xlsx.File) (*xlsx.Sheet, error) {
		if i < 0 || i >= len(file.Sheets) {
			return nil, sheetNotFoundErr
		}

		return file.Sheets[i], nil
	}
}

func OnlyFirstSheet() SheetGetter {
	return WithIndex(0)
}