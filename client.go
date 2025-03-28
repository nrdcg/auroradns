package auroradns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://api.auroradns.eu"

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

// ErrorResponse A representation of an API error message.
// Deprecated: use ResponseError instead.
type ErrorResponse = ResponseError

// ResponseError A representation of an API error message.
type ResponseError struct {
	ErrorCode string `json:"error"`
	Message   string `json:"errormsg"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%s - %s", e.ErrorCode, e.Message)
}

// Option Type of a client option.
type Option func(*Client) error

// Client The API client.
type Client struct {
	baseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

// NewClient Creates a new client.
func NewClient(httpClient *http.Client, opts ...Option) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	client := &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	req.Header.Set(contentTypeHeader, contentTypeJSON)

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if err = checkResponse(resp); err != nil {
		return resp, err
	}

	if v == nil {
		return resp, nil
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("failed to read body: %w", err)
	}

	if err = json.Unmarshal(raw, v); err != nil {
		return resp, fmt.Errorf("unmarshaling %T error: %w: %s", v, err, string(raw))
	}

	return resp, nil
}

func checkResponse(resp *http.Response) error {
	if c := resp.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err == nil && data != nil {
		errorResponse := new(ResponseError)

		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			return fmt.Errorf("unmarshaling ErrorResponse error: %w: %s", err, string(data))
		}

		return errorResponse
	}

	return fmt.Errorf("status code: %d %s", resp.StatusCode, resp.Status)
}

// WithBaseURL Allows to define a custom base URL.
func WithBaseURL(rawBaseURL string) func(*Client) error {
	return func(client *Client) error {
		if rawBaseURL == "" {
			return nil
		}

		baseURL, err := url.Parse(rawBaseURL)
		if err != nil {
			return err
		}

		client.baseURL = baseURL

		return nil
	}
}
