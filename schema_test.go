package instapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func newClient() *Client {
	return New(
		Endpoint("http://127.0.0.1:8282/v1/"),
		HTTPClient(http.DefaultClient),
		APIKey("DEADBEEFDEADBEEFDEADBEEFDEADBEEF"),
	)
}

func TestGetSchema(t *testing.T) {
	schema, err := newClient().GetSchema(context.Background(), "companies")

	require.NoError(t, err)

	spew.Dump(schema)
}

func TestDetectSchemaForFile(t *testing.T) {
	schema, err := newClient().DetectSchemaFromFile(context.Background(), "test", "testdata/companies.csv")

	require.NoError(t, err)

	spew.Dump(schema)
}
