package sec4dev

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	sdkVersion     = "1.0.0"
	connectTimeout = 10 * time.Second
	readTimeout    = 30 * time.Second
)

type rateLimitHeaders struct {
	limit        int
	remaining    int
	resetSeconds int
}

func parseRateLimit(h http.Header) rateLimitHeaders {
	getInt := func(key string) int {
		v := h.Get(key)
		if v == "" {
			return 0
		}
		n, _ := strconv.Atoi(v)
		return n
	}
	return rateLimitHeaders{
		limit:        getInt("X-RateLimit-Limit"),
		remaining:    getInt("X-RateLimit-Remaining"),
		resetSeconds: getInt("X-RateLimit-Reset"),
	}
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}) (statusCode int, out []byte, header http.Header, err error) {
	var reqBody io.Reader
	if body != nil {
		b, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return 0, nil, nil, marshalErr
		}
		reqBody = bytes.NewReader(b)
	}
	req, reqErr := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reqBody)
	if reqErr != nil {
		return 0, nil, nil, reqErr
	}
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sec4dev-go/"+sdkVersion)

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: connectTimeout + readTimeout,
		}
	}
	resp, doErr := client.Do(req)
	if doErr != nil {
		return 0, nil, nil, doErr
	}
	defer resp.Body.Close()
	out, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return resp.StatusCode, nil, nil, readErr
	}
	return resp.StatusCode, out, resp.Header.Clone(), nil
}

func isRetryable(statusCode int, netErr bool) bool {
	if netErr {
		return true
	}
	if statusCode == 429 || statusCode == 500 || statusCode == 502 || statusCode == 503 || statusCode == 504 {
		return true
	}
	return false
}

func parseErrorBody(body []byte) (message string, parsed interface{}) {
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return "Unknown error", body
	}
	if d, ok := m["detail"]; ok {
		if s, ok := d.(string); ok {
			return s, m
		}
		return "Unknown error", m
	}
	return "Unknown error", m
}

// postWithRetry performs POST with retries and returns response body and rate limit info.
func (c *Client) postWithRetry(ctx context.Context, path string, body interface{}, onRateLimit func(RateLimitInfo)) ([]byte, RateLimitInfo, error) {
	var lastErr error
	var lastStatus int
	var lastBody []byte
	var lastHeader http.Header
	rl := RateLimitInfo{}

	for attempt := 0; attempt <= c.Retries; attempt++ {
		status, out, header, err := c.do(ctx, "POST", path, body)
		if err != nil {
			lastErr = err
			if attempt < c.Retries {
				delay := time.Duration(c.RetryDelayMs)*time.Millisecond*time.Duration(1<<attempt) + time.Duration(rand.Intn(101))*time.Millisecond
				select {
				case <-ctx.Done():
					return nil, rl, ctx.Err()
				case <-time.After(delay):
					continue
				}
			}
			return nil, rl, err
		}

		rh := parseRateLimit(header)
		rl = RateLimitInfo{Limit: rh.limit, Remaining: rh.remaining, ResetSeconds: rh.resetSeconds}
		if onRateLimit != nil {
			onRateLimit(rl)
		}

		if status == 429 {
			retryAfter := 60
			if v := header.Get("Retry-After"); v != "" {
				if n, _ := strconv.Atoi(v); n > 0 {
					retryAfter = n
				}
			}
			if attempt < c.Retries {
				select {
				case <-ctx.Done():
					return nil, rl, ctx.Err()
				case <-time.After(time.Duration(retryAfter) * time.Second):
					continue
				}
			}
			msg, parsed := parseErrorBody(out)
			return nil, rl, errFromStatus(429, msg, parsed, retryAfter, rh.limit, rh.remaining)
		}

		if status >= 400 {
			msg, parsed := parseErrorBody(out)
			apiErr := errFromStatus(status, msg, parsed, 0, rh.limit, rh.remaining)
			if !isRetryable(status, false) {
				return nil, rl, apiErr
			}
			lastErr = apiErr
			lastStatus = status
			lastBody = out
			lastHeader = header
			if attempt < c.Retries {
				delay := time.Duration(c.RetryDelayMs)*time.Millisecond*time.Duration(1<<attempt) + time.Duration(rand.Intn(101))*time.Millisecond
				select {
				case <-ctx.Done():
					return nil, rl, ctx.Err()
				case <-time.After(delay):
					continue
				}
			}
			return nil, rl, apiErr
		}

		return out, rl, nil
	}

	if lastErr != nil && lastStatus >= 400 && lastHeader != nil {
		msg, parsed := parseErrorBody(lastBody)
		rh := parseRateLimit(lastHeader)
		return nil, rl, errFromStatus(lastStatus, msg, parsed, 0, rh.limit, rh.remaining)
	}
	if lastErr != nil {
		return nil, rl, lastErr
	}
	return nil, rl, baseError("Request failed after retries", 0, nil)
}
