package instapi

import (
	"net/url"
	"strconv"
)

// RequestOption represents a API request option.
type RequestOption func(*url.Values)

func Name(name string) RequestOption {
	return func(u *url.Values) {
		u.Set("name", name)
	}
}

func Limit(limit int) RequestOption {
	return func(u *url.Values) {
		u.Set("limit", strconv.Itoa(limit))
	}
}

func Skip(skip int) RequestOption {
	return func(u *url.Values) {
		u.Set("skip", strconv.Itoa(skip))
	}
}

func Table(table string) RequestOption {
	return func(u *url.Values) {
		u.Set("table", table)
	}
}

func FromCell(cell string) RequestOption {
	return func(u *url.Values) {
		u.Set("fromCell", cell)
	}
}

func Offset(offset string) RequestOption {
	return func(u *url.Values) {
		u.Set("offset", offset)
	}
}
