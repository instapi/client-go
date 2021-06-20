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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/instapi/client-go/internal/csvutil"
	"github.com/instapi/client-go/types"

	"github.com/tomnomnom/linkheader"
)

// Instapi client constants.
const (
	DefaultEndpoint = "https://api.instapi.com/v1/"
)

// Client related errors.
var (
	ErrNotFound        = errors.New("resource not found")
	ErrUnsupportedType = errors.New("unsupported type")
	ErrForbidden       = errors.New("forbidden")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrStatus          = errors.New("unexpected HTTP status")
)

// Client represents a client implementation.
type Client struct {
	doer      Doer
	debugFunc func(*http.Request, *http.Response, Debug)
	endpoint  string
	token     string
}

// Doer defines the HTTP Do() interface.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientOption represents a client option.
type ClientOption func(*Client)

// HTTPClient option.
func HTTPClient(doer Doer) ClientOption {
	return func(c *Client) {
		c.doer = doer
	}
}

// Endpoint option.
func Endpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// Token option.
func Token(token string) ClientOption {
	return func(c *Client) {
		c.token = "Bearer " + token
	}
}

// DebugFunc option.
func DebugFunc(f func(*http.Request, *http.Response, Debug)) ClientOption {
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

// Default returns a default client configured from environment variables.
func Default() *Client {
	return New(Endpoint(os.Getenv("API_ENDPOINT")), Token(os.Getenv("TOKEN")))
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.token)

	return req, nil
}

// Debug information returned by a debug function.
type Debug struct {
	Payload  []byte
	Duration time.Duration
}

func (c *Client) doRequest(ctx context.Context, method, contentType, endpoint string, statusCode int, src, dst interface{}, options ...RequestOption) (*http.Response, []byte, error) {
	var (
		r       io.Reader
		payload []byte
	)

	if src != nil {
		if v, ok := src.(io.Reader); ok {
			// Optimize sending large CSV payloads with a record limit
			if contentType == types.CSV {
				limit := 0
				headers := true // Headers are assumed by default

				for _, option := range options {
					switch option.param {
					case "limit":
						limit = option.value.(int)

					case "headers":
						headers = option.value.(bool)
					}
				}

				if limit > 0 {
					if headers {
						limit++
					}

					r = csvutil.NewLineLimitReader(v, limit)
				}
			}

			if r == nil {
				r = v
			}
		} else {
			var err error
			payload, err = json.Marshal(src)

			if err != nil {
				return nil, nil, err
			}

			r = bytes.NewReader(payload)
		}
	}

	req, err := c.newRequest(ctx, method, endpoint, r)

	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Accept", contentType)
	req.Header.Add("Content-Type", contentType)

	nilDst := dst == nil

	if nilDst {
		switch method {
		case http.MethodPatch, http.MethodPost, http.MethodPut:
			req.Header.Add("No-Response-Body", "1")
		}
	}

	q := req.URL.Query()

	for _, option := range options {
		option.fn(&q)
	}

	req.URL.RawQuery = q.Encode()
	start := time.Now()
	resp, err := c.doer.Do(req)

	if err != nil {
		return nil, nil, err
	}

	d := time.Since(start)

	defer resp.Body.Close() // nolint: errcheck

	if c.debugFunc != nil {
		c.debugFunc(req, resp, Debug{Payload: payload, Duration: d})
	}

	// Early exit for successful HTTP status code and nil destination
	if nilDst &&
		(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) &&
		(statusCode == http.StatusOK || statusCode == http.StatusNoContent) {
		return resp, nil, nil
	}

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, nil, err
	}

	if statusCode > 0 &&
		resp.StatusCode != statusCode ||
		(resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices) {
		err := decodeAPIError(resp.StatusCode, b)

		if err != nil {
			return nil, nil, err
		}

		switch resp.StatusCode {
		case http.StatusForbidden:
			return nil, nil, fmt.Errorf("%w: %s %s", ErrForbidden, method, endpoint)
		case http.StatusNotFound:
			return nil, nil, fmt.Errorf("%w: %s", ErrNotFound, endpoint)
		case http.StatusUnauthorized:
			return nil, nil, fmt.Errorf("%w: %s %s", ErrUnauthorized, method, endpoint)
		default:
			return nil, nil, fmt.Errorf("%w: expected %d, got %d", ErrStatus, statusCode, resp.StatusCode)
		}
	}

	if dst != nil {
		return resp, b, json.Unmarshal(b, &dst)
	}

	return resp, b, nil
}

// Error represents a client error.
type Error struct {
	StatusCode int
	Err        string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func decodeAPIError(statusCode int, b []byte) error {
	if len(b) == 0 || b[0] != '{' {
		return nil
	}

	var e Error
	err := json.Unmarshal(b, &e)

	if err != nil {
		return err
	}

	if e.Err == "" {
		return nil
	}

	e.StatusCode = statusCode

	return e
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
