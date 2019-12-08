package client

import (
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func newClient() *Client {
	return New(
		WithEndpoint("http://127.0.0.1:8282/v1/"),
		WithHTTPClient(http.DefaultClient),
		WithAPIKey("DEADBEEFDEADBEEFDEADBEEFDEADBEEF"),
	)
}

func TestGetSchema(t *testing.T) {
	schema, err := newClient().GetSchema("companies")

	require.NoError(t, err)

	spew.Dump(schema)
}

func TestDetectSchemaForFile(t *testing.T) {
	schema, err := newClient().DetectSchemaForFile("test", "testdata/companies.csv")

	require.NoError(t, err)

	spew.Dump(schema)
}
