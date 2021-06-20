package instapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/instapi/client-go/types"
)

// Query performs a SQL query.
func (c *Client) Query(ctx context.Context, query string, dst interface{}) error {
	_, _, err := c.doRequest(
		ctx,
		http.MethodPost,
		types.SQL,
		c.endpoint+"query",
		http.StatusOK,
		strings.NewReader(query),
		dst,
	)

	return err
}
