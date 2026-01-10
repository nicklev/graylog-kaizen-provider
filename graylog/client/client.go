package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Graylog API client
type Client struct {
	BaseURL      string
	Username     string
	Password     string
	HTTPClient   *http.Client
	XRequestedBy string
	APIVersion   string
}

// NewClient creates a new Graylog client
func NewClient(baseURL, username, password *string) (*Client, error) {
	if baseURL == nil || *baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if username == nil || *username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if password == nil || *password == "" {
		return nil, fmt.Errorf("password is required")
	}

	return &Client{
		BaseURL:  *baseURL,
		Username: *username,
		Password: *password,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
		XRequestedBy: "terraform-provider-graylog",
		APIVersion:   "v3", // Default to v3, can be overridden
	}, nil
}

// SetXRequestedBy sets the X-Requested-By header value
func (c *Client) SetXRequestedBy(value string) {
	c.XRequestedBy = value
}

// SetAPIVersion sets the API version
func (c *Client) SetAPIVersion(version string) {
	c.APIVersion = version
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s/api/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.XRequestedBy != "" {
		req.Header.Set("X-Requested-By", c.XRequestedBy)
	}

	// Set basic authentication
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// Get performs a GET request
func (c *Client) Get(endpoint string, result interface{}) error {
	resp, err := c.doRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Post performs a POST request
func (c *Client) Post(endpoint string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Put performs a PUT request
func (c *Client) Put(endpoint string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(http.MethodPut, endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Delete performs a DELETE request
func (c *Client) Delete(endpoint string) error {
	resp, err := c.doRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
