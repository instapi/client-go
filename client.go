package instapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-querystring/query"
	"github.com/instapi/client-go/schema"
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

// DetectSchemaForFile infers the schema for the given file.
func (c *Client) DetectSchemaForFile(options *DetectOptions, filename string) (*schema.Schema, error) {
	ext := filepath.Ext(filename)
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close() // nolint: errcheck

	return c.DetectSchema(options, f, ext)
}

// DetectOptions represents schema detection options.
type DetectOptions struct {
	Name     string `url:"name"`
	Table    string `url:"table,omitempty"`
	FromCell string `url:"fromCell,omitempty"`
	Limit    int    `url:"limit,omitempty"`
}

// DetectSchema attempts to detect the schema for the give reader.
func (c *Client) DetectSchema(options *DetectOptions, r io.Reader, ext string) (*schema.Schema, error) {
	contentType, err := schema.ExtContentType(ext)

	if err != nil {
		return nil, err
	}

	query, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(http.MethodPost, c.endpoint+"schemas/detect", r)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	req.URL.RawQuery = query.Encode()
	resp, err := c.doer.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Body", string(b))

		return nil, fmt.Errorf("%d - %w", resp.StatusCode, ErrStatus)
	}

	var schema *schema.Schema

	return schema, json.NewDecoder(resp.Body).Decode(&schema)
}
