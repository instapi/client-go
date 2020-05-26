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
	"strings"

	"github.com/instapi/client-go/types"
	"github.com/tomnomnom/linkheader"
)

// Instapi client constants.
const (
	DefaultEndpoint = "https://api.instapi.com/v1/"
)

// Client related errors.
var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrForbidden       = errors.New("forbidden")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrStatus          = errors.New("unexpected HTTP status")
)

// Client represents a client implementation.
type Client struct {
	doer      Doer
	debugFunc func(*http.Response)
	endpoint  string
	apiKey    string
}

// Doer defines the HTTP Do() interface.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientOption represents a client option.
type ClientOption func(*Client)

func HTTPClient(doer Doer) ClientOption {
	return func(c *Client) {
		c.doer = doer
	}
}

func Endpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

func APIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = "Key " + key
	}
}

func DebugFunc(f func(*http.Response)) ClientOption {
	return func(c *Client) {
		c.debugFunc = f
	}
}

// New initializes a new client instance.
func New(options ...ClientOption) *Client {
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

func (c *Client) doRequest(ctx context.Context, method, contentType, endpoint string, statusCode int, src, dst interface{}, options ...RequestOption) (*http.Response, error) {
	var r io.Reader

	if src != nil {
		if v, ok := src.(io.Reader); ok {
			r = v
		} else {
			buf := &bytes.Buffer{}
			err := json.NewEncoder(buf).Encode(src)

			if err != nil {
				return nil, err
			}

			r = buf
		}
	}

	req, err := c.newRequest(ctx, method, endpoint, r)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", contentType)
	req.Header.Add("Content-Type", contentType)

	if dst == nil {
		switch method {
		case http.MethodPatch, http.MethodPost, http.MethodPut:
			req.Header.Add("No-Response-Body", "1")
		}
	}

	q := req.URL.Query()

	for _, option := range options {
		option(&q)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := c.doer.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint: errcheck

	if c.debugFunc != nil {
		c.debugFunc(resp)
	}

	if statusCode > 0 &&
		resp.StatusCode != statusCode || (resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices) {
		switch resp.StatusCode {
		case http.StatusForbidden:
			return nil, fmt.Errorf("%s %s - %w", method, endpoint, ErrForbidden)
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("%s %s - %w", method, endpoint, ErrUnauthorized)
		default:
			return nil, fmt.Errorf("expected %d, got %d - %w", statusCode, resp.StatusCode, ErrStatus)
		}
	}

	if dst != nil {
		return resp, json.NewDecoder(resp.Body).Decode(dst)
	}

	return resp, nil
}

func nextLink(resp *http.Response) (string, error) {
	for _, v := range linkheader.Parse(strings.TrimPrefix(resp.Header.Get("link"), "Link:")) {
		if v.Rel != "next" {
			continue
		}

		u, err := url.Parse(v.URL)

		if err != nil {
			return "", err
		}

		return u.Query().Get("offset"), nil
	}

	return "", nil
}

func getContentType(filename string) (string, error) {
	return types.TypeFromExt(filepath.Ext(filename))
}

func escape(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}
