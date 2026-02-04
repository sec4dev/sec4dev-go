package sec4dev

import (
	"net/http"
	"strings"
)

const defaultBaseURL = "https://api.sec4.dev/api/v1"

// Client is the Sec4Dev API client.
type Client struct {
	APIKey       string
	BaseURL      string
	HTTPClient   *http.Client
	Retries      int
	RetryDelayMs int
	onRateLimit  func(RateLimitInfo)
	rateLimit    RateLimitInfo
}

// ClientOption configures the client.
type ClientOption func(*Client)

// WithBaseURL sets the API base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = strings.TrimSuffix(url, "/")
	}
}

// WithTimeout is not used directly; use a custom HTTPClient with timeout.
// WithRetries sets the number of retries.
func WithRetries(n int) ClientOption {
	return func(c *Client) {
		c.Retries = n
	}
}

// WithRetryDelay sets the base retry delay in milliseconds.
func WithRetryDelay(ms int) ClientOption {
	return func(c *Client) {
		c.RetryDelayMs = ms
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = hc
	}
}

// WithRateLimitCallback sets a callback for rate limit updates.
func WithRateLimitCallback(fn func(RateLimitInfo)) ClientOption {
	return func(c *Client) {
		c.onRateLimit = fn
	}
}

// NewClient creates a new Sec4Dev client. API key must start with "sec4_".
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	key := strings.TrimSpace(apiKey)
	if key == "" || !strings.HasPrefix(key, "sec4_") {
		return nil, &ValidationError{baseError("API key must start with sec4_", 422, nil)}
	}
	c := &Client{
		APIKey:       key,
		BaseURL:      defaultBaseURL,
		Retries:      3,
		RetryDelayMs: 1000,
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// RateLimit returns the last rate limit info.
func (c *Client) RateLimit() RateLimitInfo {
	return c.rateLimit
}

// Email returns the email service.
func (c *Client) Email() *EmailService {
	return &EmailService{client: c}
}

// IP returns the IP service.
func (c *Client) IP() *IPService {
	return &IPService{client: c}
}
