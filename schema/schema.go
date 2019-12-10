package schema

import (
	"errors"
	"fmt"
)

const (
	mimeXLS  = "application/vnd.ms-excel"
	mimeXLSX = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

// Package errors.
var (
	ErrUnsupportedExtension = errors.New("unsupported extension")
)

var (
	contentTypes = map[string]string{
		".csv":  "text/csv",
		".ods":  "application/vnd.oasis.opendocument.spreadsheet",
		".xls":  mimeXLS,
		".xla":  mimeXLS,
		".xlsx": mimeXLSX,
		".xlsm": mimeXLSX,
		".xlam": mimeXLSX,
		".xlsb": "application/vnd.ms-excel.sheet.binary.macroEnabled.12",
	}
)

// Schema represents a schema.
type Schema struct {
	Name       string    `json:"name"`
	PrimaryKey []string  `json:"primaryKey,omitempty"`
	Count      int       `json:"count,omitempty"`
	Settings   *Settings `json:"settings,omitempty"`
	Fields     []Field   `json:"fields"`
}

// Settings represents schema settings.
type Settings struct {
	QuoteValues *string `json:"quoteValues,omitempty"`
}

// Field represents a schema field.
type Field struct {
	Name     string `json:"name"`
	Mapping  string `json:"mapping,omitempty"`
	Type     string `json:"type"`
	Required bool   `json:"required,omitempty"`
}

// DetectOptions represents schema detection options.
type DetectOptions struct {
	Name     string `url:"name"`
	Table    string `url:"table,omitempty"`
	FromCell string `url:"fromCell,omitempty"`
	Limit    int    `url:"limit,omitempty"`
}

// ContentType returns the content type for the given file extension.
func ContentType(ext string) (string, error) {
	v, exists := contentTypes[ext]

	if !exists {
		return "", fmt.Errorf("extension: %s - %w", ext, ErrUnsupportedExtension)
	}

	return v, nil
}
