package xlsx2csv

import (
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	sheetNotFoundErr = errors.New("requested sheet not found")
)

type SheetGetter func(file *excelize.File) (string, error)

func WithName(name string) SheetGetter {
	return func(file *excelize.File) (string, error) {
		return name, nil
	}
}

func WithIndex(i int) SheetGetter {
	return func(file *excelize.File) (string, error) {
		if i < 0 || i > file.SheetCount {
			return "", sheetNotFoundErr
		}
		return file.GetSheetName(i), nil
	}
}

func OnlyFirstSheet() SheetGetter {
	return WithIndex(1)
}