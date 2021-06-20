package instapi

import (
	"context"
	"net/http"
	"net/url"

	"github.com/instapi/client-go/record"
	"github.com/instapi/client-go/schema"
	"github.com/instapi/client-go/types"
)

// DetectSheetSchemas attempts to detect the schema for the given Google Sheets sheet.
func (c *Client) DetectSheetSchemas(ctx context.Context, name, sheetID, rng string, options ...RequestOption) ([]*schema.Schema, error) {
	var s []*schema.Schema
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"sheets/detect",
		http.StatusOK,
		nil,
		&s,
		append(options, Name(name), ExternalID(sheetID), Range(rng))...,
	)

	return s, err
}

// ImportSheet imports the given Google Sheet.
func (c *Client) ImportSheet(ctx context.Context, account, schema, sheetID, rng string, options ...RequestOption) (int, error) {
	var resp record.Batch
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.JSON,
		c.endpoint+"accounts/"+url.PathEscape(account)+"/schemas/"+url.PathEscape(schema)+"/records",
		http.StatusAccepted,
		nil,
		&resp,
		append(options, ExternalID(sheetID), Range(rng))...,
	)

	return resp.Count, err
}
