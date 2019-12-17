package instapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/instapi/client-go/types"
)

// Instapi client constants.
const (
	DefaultEndpoint = "https://api.instapi.com/v1/"
)

// Client related errors.
var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrStatus          = errors.New("unexpected HTTP status")
)

// Client represents a client implementation.
type Client struct {
	doer     Doer
	endpoint string
	apiKey   string
}

// Options represents shared options.
type Options struct {
	Name       string `url:"name,omitempty"`
	PrimaryKey string `url:"primaryKey,omitempty"`
	SortKey    string `url:"sortKey,omitempty"`
	Table      string `url:"table,omitempty"`
	Sheet      string `url:"sheet,omitempty"`
	FromCell   string `url:"fromCell,omitempty"`
	ToCell     string `url:"toCell,omitempty"`
	Skip       int    `url:"skip,omitempty"`
	Limit      int    `url:"limit,omitempty"`
	Headers    bool   `url:"headers,omitempty"`
}

// Doer defines the HTTP Do() interface.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Option represents a client option.
type Option func(*Client)

func WithHTTPClient(doer Doer) Option {
	return func(c *Client) {
		c.doer = doer
	}
}

func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = "Key " + key
	}
}

// New initializes a new client instance.
func New(options ...Option) *Client {
	c := &Client{}

	for _, option := range options {
		option(c)
	}

	if c.endpoint == "" {
		c.endpoint = DefaultEndpoint
	}

	if c.doer == nil {
		c.doer = http.DefaultClient
	}

	return c
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.apiKey)

	return req, nil
}

func (c *Client) doRequest(ctx context.Context, method, contentType, endpoint string, statusCode int, query url.Values, src, dst interface{}) error {
	var r io.Reader

	if src != nil {
		if v, ok := src.(io.Reader); ok {
			r = v
		} else {
			buf := &bytes.Buffer{}
			err := json.NewEncoder(buf).Encode(src)

			if err != nil {
				return err
			}

			r = buf
		}
	}

	req, err := c.newRequest(ctx, method, endpoint, r)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)
	req.URL.RawQuery = query.Encode()
	resp, err := c.doer.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != statusCode {
		return fmt.Errorf("expected %d, got %d - %w", statusCode, resp.StatusCode, ErrStatus)
	}

	if dst != nil {
		return json.NewDecoder(resp.Body).Decode(dst)
	}

	return nil
}

func getContentType(filename string) (string, error) {
	return types.TypeFromExt(filepath.Ext(filename))
}
