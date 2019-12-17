package instapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/google/go-querystring/query"
)

// CreateRecords makes a create records request.
func (c *Client) CreateRecords(ctx context.Context, options *Options, schema, contentType string, r io.Reader) (int, error) {
	query, err := query.Values(options)

	if err != nil {
		return 0, err
	}

	var resp struct {
		Count int `json:"count"`
	}

	return resp.Count, c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records",
		http.StatusCreated,
		query,
		r,
		&resp,
	)
}

// CreateRecordsFromFile makes a create records request for the given file.
func (c *Client) CreateRecordsFromFile(ctx context.Context, options *Options, schema, filename string) (int, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return 0, err
	}

	f, err := os.Open(filename)

	if err != nil {
		return 0, err
	}

	defer f.Close() // nolint: errcheck

	return c.CreateRecords(ctx, options, schema, contentType, f)
}
