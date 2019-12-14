package instapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/go-querystring/query"
	"github.com/instapi/client-go/schema"
)

// Instapi client constants.
const (
	DefaultEndpoint    = "https://api.instapi.com/v1/"
	defaultContentType = "application/json"
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

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.apiKey)

	return req, nil
}

func (c *Client) doRequest(method, endpoint string, statusCode int, contentType *string, query url.Values, src, dst interface{}) error {
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

	req, err := c.newRequest(method, endpoint, r)

	if err != nil {
		return err
	}

	if contentType != nil {
		req.Header.Add("Content-Type", *contentType)
	} else {
		req.Header.Add("Content-Type", defaultContentType)
	}

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

// GetSchema gets the given schema.
func (c *Client) GetSchema(name string) (*schema.Schema, error) {
	req, err := c.newRequest(http.MethodGet, c.endpoint+"schemas/"+name, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.doer.Do(req)

	if err != nil {
		return nil, err
	}

	var schema *schema.Schema

	return schema, json.NewDecoder(resp.Body).Decode(&schema)
}

// DetectSchemaForFile attempts to detect the schema for the given file.
func (c *Client) DetectSchemaForFile(options *schema.DetectOptions, filename string) (*schema.Schema, error) {
	ext := filepath.Ext(filename)
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close() // nolint: errcheck

	return c.DetectSchema(options, f, ext)
}

// DetectSchema attempts to detect the schema for the given reader.
func (c *Client) DetectSchema(options *schema.DetectOptions, r io.Reader, ext string) (*schema.Schema, error) {
	contentType, err := schema.ContentType(ext)

	if err != nil {
		return nil, err
	}

	query, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	var s *schema.Schema

	return s, c.doRequest(
		http.MethodPost,
		c.endpoint+"schemas/detect",
		http.StatusOK,
		stringPtr(contentType),
		query,
		r,
		&s,
	)
}

// CreateSchema creates a new schema.
func (c *Client) CreateSchema(s *schema.Schema) error {
	return c.doRequest(
		http.MethodPost,
		c.endpoint+"schemas",
		http.StatusCreated,
		nil,
		nil,
		s,
		nil,
	)
}

// DetectAndCreateSchemaForFile attempts to detect and create
// the schema for the given file.
func (c *Client) DetectAndCreateSchemaForFile(options *schema.DetectOptions, filename string) (*schema.Schema, error) {
	s, err := c.DetectSchemaForFile(options, filename)

	if err != nil {
		return nil, err
	}

	err = c.CreateSchema(s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DetectAndCreateSchema attempts to detect and create
// the schema for the given reader.
func (c *Client) DetectAndCreateSchema(options *schema.DetectOptions, r io.Reader, ext string) (*schema.Schema, error) {
	s, err := c.DetectSchema(options, r, ext)

	if err != nil {
		return nil, err
	}

	err = c.CreateSchema(s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteSchema deletes a schema.
func (c *Client) DeleteSchema(name string) error {
	return c.doRequest(
		http.MethodDelete,
		c.endpoint+"schemas/"+url.PathEscape(name),
		http.StatusAccepted,
		nil,
		nil,
		nil,
		nil,
	)
}

func stringPtr(s string) *string { return &s }
