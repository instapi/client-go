package types

import (
	"errors"
	"fmt"
)

// Package errors.
var (
	ErrUnsupportedExtension = errors.New("unsupported extension")
)

// Constants for types package.
const (
	mimeXLS  = "application/vnd.ms-excel"
	mimeXLSX = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	CSV      = "text/csv"
	JSON     = "application/json"
	ODS      = "application/vnd.oasis.opendocument.spreadsheet"
	SQLite   = "application/x-sqlite3"
	XLA      = mimeXLS
	XLS      = mimeXLS
	XLAM     = mimeXLSX
	XLSB     = "application/vnd.ms-excel.sheet.binary.macroEnabled.12"
	XLSM     = mimeXLSX
	XLSX     = mimeXLSX
)

var contentTypes = map[string]string{
	".csv":  CSV,
	".db":   SQLite,
	".json": JSON,
	".ods":  ODS,
	".xls":  XLS,
	".xla":  XLA,
	".xlsx": XLSX,
	".xlsm": XLSM,
	".xlam": XLAM,
	".xlsb": XLSB,
}

// TypeFromExt returns the content type for the given file extension.
func TypeFromExt(ext string) (string, error) {
	v, exists := contentTypes[ext]

	if !exists {
		return "", fmt.Errorf("extension: %s - %w", ext, ErrUnsupportedExtension)
	}

	return v, nil
}
