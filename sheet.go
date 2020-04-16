package xlsx2csv

import (
	"errors"

	"github.com/tealeg/xlsx/v2"
)

var (
	sheetNotFoundErr = errors.New("requested sheet not found")
)

type SheetSelector func(file *xlsx.File) (*xlsx.Sheet, error)

func SheetByName(name string) SheetSelector {
	return func(file *xlsx.File) (*xlsx.Sheet, error) {
		sheet, ok := file.Sheet[name]
		if !ok {
			return nil, sheetNotFoundErr
		}

		return sheet, nil
	}
}

func SheetByIndex(i int) SheetSelector {
	return func(file *xlsx.File) (*xlsx.Sheet, error) {
		if i < 0 || i >= len(file.Sheets) {
			return nil, sheetNotFoundErr
		}

		return file.Sheets[i], nil
	}
}

func FirstSheet() SheetSelector {
	return SheetByIndex(0)
}