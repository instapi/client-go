package instapi

import (
	"fmt"
	"net/url"
	"strconv"
)

// RequestOption represents a API request option.
type RequestOption struct {
	fn    func(*url.Values)
	param string
	value interface{}
}

// Schema sets the schema parameter.
func Schema(schema string) RequestOption {
	return Param("schema", schema)
}

// Name sets the name parameter.
func Name(name string) RequestOption {
	return Param("name", name)
}

// Limit sets the limit parameter.
func Limit(limit int) RequestOption {
	return Param("limit", limit)
}

// Skip sets the skip parameter.
func Skip(skip int) RequestOption {
	return Param("skip", skip)
}

// Table sets the table parameter.
func Table(table string) RequestOption {
	return Param("table", table)
}

// CellFrom sets the from cell parameter.
func CellFrom(cell string) RequestOption {
	return Param("cellFrom", cell)
}

// CellTo sets the to cell parameter.
func CellTo(cell string) RequestOption {
	return Param("cellTo", cell)
}

// ExternalID sets the to external ID parameter.
func ExternalID(externalID string) RequestOption {
	return Param("externalId", externalID)
}

// Range sets the to range parameter.
func Range(rng string) RequestOption {
	return Param("range", rng)
}

// Offset sets the from offset parameter.
func Offset(offset string) RequestOption {
	return Param("offset", offset)
}

// Headers sets the from headers parameter.
func Headers(headers bool) RequestOption {
	return Param("headers", headers)
}

// Param sets given URL parameter with the given value.
func Param(k string, v interface{}) RequestOption {
	value := ""

	switch v := v.(type) {
	case int:
		value = strconv.Itoa(v)

	case bool:
		if v {
			value = "true"
		} else {
			value = "false"
		}

	case string:
		value = v

	default:
		s, ok := v.(fmt.Stringer)

		if ok {
			value = s.String()
		} else {
			value = fmt.Sprintf("%v", v)
		}
	}

	return RequestOption{
		fn: func(u *url.Values) {
			u.Set(k, value)
		},
		param: k,
		value: v,
	}
}
