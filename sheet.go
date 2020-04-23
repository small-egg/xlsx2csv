package xlsx2csv

import (
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	sheetNotFoundErr = errors.New("requested sheet not found")
)

type SheetSelector func(file *excelize.File) (string, error)

func SheetByName(name string) SheetSelector {
	return func(file *excelize.File) (string, error) {
		return name, nil
	}
}

func SheetByIndex(i int) SheetSelector {
	return func(file *excelize.File) (string, error) {
		if i < 0 || i > file.SheetCount {
			return "", sheetNotFoundErr
		}

		return file.GetSheetList()[i], nil
	}
}

func FirstSheet() SheetSelector {
	return SheetByIndex(0)
}
