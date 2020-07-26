package instapi

import (
	"net/url"
	"strconv"
)

// RequestOption represents a API request option.
type RequestOption func(*url.Values)

// Name sets the name parameter.
func Name(name string) RequestOption {
	return Param("name", name)
}

// Limit sets the limit parameter.
func Limit(limit int) RequestOption {
	return Param("limit", strconv.Itoa(limit))
}

// Skip sets the skip parameter.
func Skip(skip int) RequestOption {
	return Param("skip", strconv.Itoa(skip))
}

// Table sets the table parameter.
func Table(table string) RequestOption {
	return Param("table", table)
}

// FromCell sets the from cell parameter.
func FromCell(cell string) RequestOption {
	return Param("fromCell", cell)
}

// Offset sets the from offset parameter.
func Offset(offset string) RequestOption {
	return Param("offset", offset)
}

// Param sets given URL parameter with the given value.
func Param(k, v string) RequestOption {
	return func(u *url.Values) {
		u.Set(k, v)
	}
}
