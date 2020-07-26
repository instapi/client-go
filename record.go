package instapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/instapi/client-go/types"
)

// GetRecords gets schema records.
func (c *Client) GetRecords(ctx context.Context, schema string, dst interface{}, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records",
		http.StatusOK,
		nil,
		dst,
		options...,
	)

	return err
}

// GetRecord gets a record.
func (c *Client) GetRecord(ctx context.Context, schema, id string, dst interface{}, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records/"+escape(id),
		http.StatusOK,
		nil,
		dst,
		options...,
	)

	return err
}

// CreateRecord makes a create record request.
func (c *Client) CreateRecord(ctx context.Context, schema string, src, dst interface{}, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records",
		http.StatusCreated,
		src,
		dst,
		options...,
	)

	return err
}

// CreateRecords makes a create records request.
func (c *Client) CreateRecords(ctx context.Context, schema, contentType string, r io.Reader, options ...RequestOption) (int, error) {
	var resp struct {
		Count int `json:"count"`
	}
	_, err := c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records",
		http.StatusAccepted,
		r,
		&resp,
		append(options, func(v *url.Values) {
			v.Set("batch", "true")
		})...,
	)

	return resp.Count, err
}

// CreateRecordsFromFile makes a create records request for the given file.
func (c *Client) CreateRecordsFromFile(ctx context.Context, schema, filename string, options ...RequestOption) (int, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return 0, err
	}

	f, err := os.Open(filename)

	if err != nil {
		return 0, err
	}

	defer f.Close() // nolint: errcheck

	return c.CreateRecords(ctx, schema, contentType, f, options...)
}

// UpdateRecord updates a record.
func (c *Client) UpdateRecord(ctx context.Context, schema, id string, src interface{}, dst interface{}, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodPut,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records/"+escape(id),
		http.StatusOK,
		src,
		dst,
		options...,
	)

	return err
}

// PatchRecord patches a record.
func (c *Client) PatchRecord(ctx context.Context, schema, id string, src interface{}, dst interface{}, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodPatch,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records/"+escape(id),
		http.StatusOK,
		src,
		dst,
		options...,
	)

	return err
}

// DeleteRecord deletes a record.
func (c *Client) DeleteRecord(ctx context.Context, schema, id string, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records/"+escape(id),
		http.StatusNoContent,
		nil,
		nil,
		options...,
	)

	return err
}

// DeleteRecords deletes a batch of records.
func (c *Client) DeleteRecords(ctx context.Context, schema string, IDs []string, options ...RequestOption) error {
	_, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"schemas/"+url.PathEscape(schema)+"/records",
		http.StatusNoContent,
		IDs,
		nil,
		options...,
	)

	return err
}
