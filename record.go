package instapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/instapi/client-go/record"
	"github.com/instapi/client-go/types"
)

// GetRecords gets schema records.
func (c *Client) GetRecords(ctx context.Context, account, schema string, dst interface{}, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records",
		http.StatusOK,
		nil,
		dst,
		options...,
	)

	return err
}

// GetRecord gets a record.
func (c *Client) GetRecord(ctx context.Context, account, schema, id string, dst interface{}, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodGet,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records/"+id,
		http.StatusOK,
		nil,
		dst,
		options...,
	)

	return err
}

// CreateRecord makes a create record request.
func (c *Client) CreateRecord(ctx context.Context, account, schema string, src, dst interface{}, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records",
		http.StatusCreated,
		src,
		dst,
		options...,
	)

	return err
}

// CreateRecords makes a create records request.
func (c *Client) CreateRecords(ctx context.Context, account, schema, contentType string, r io.Reader, options ...RequestOption) (int, error) {
	var resp record.Batch
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		contentType,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records",
		http.StatusAccepted,
		r,
		&resp,
		append(options, Param("batch", true))...,
	)

	return resp.Count, err
}

// CreateRecordsFromFile makes a create records request for the given file.
func (c *Client) CreateRecordsFromFile(ctx context.Context, account, schema, filename string, options ...RequestOption) (int, error) {
	contentType, err := getContentType(filename)

	if err != nil {
		return 0, err
	}

	f, err := os.Open(filename) // nolint: gosec

	if err != nil {
		return 0, err
	}

	defer f.Close() // nolint: gosec

	return c.CreateRecords(ctx, account, schema, contentType, f, options...)
}

// UpdateRecord updates a record.
func (c *Client) UpdateRecord(ctx context.Context, account, schema, id string, src interface{}, dst interface{}, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPut,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records/"+id,
		http.StatusOK,
		src,
		dst,
		options...,
	)

	return err
}

// PatchRecord patches a record.
func (c *Client) PatchRecord(ctx context.Context, account, schema, id string, src interface{}, dst interface{}, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPatch,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records/"+id,
		http.StatusOK,
		src,
		dst,
		options...,
	)

	return err
}

// DeleteRecord deletes a record.
func (c *Client) DeleteRecord(ctx context.Context, account, schema, id string, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records/"+id,
		http.StatusNoContent,
		nil,
		nil,
		options...,
	)

	return err
}

// DeleteRecords deletes a batch of records.
func (c *Client) DeleteRecords(ctx context.Context, account, schema string, ids []string, options ...RequestOption) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodDelete,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records",
		http.StatusNoContent,
		ids,
		nil,
		options...,
	)

	return err
}
