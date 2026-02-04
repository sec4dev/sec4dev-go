package sec4dev

import (
	"context"
	"encoding/json"
	"strings"
)

// EmailService provides email check operations.
type EmailService struct {
	client *Client
}

// Check checks if an email uses a disposable domain.
func (s *EmailService) Check(ctx context.Context, email string) (*EmailCheckResult, error) {
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	path := "/email/check"
	body := map[string]string{"email": strings.TrimSpace(email)}
	var onRateLimit func(RateLimitInfo)
	if s.client.onRateLimit != nil {
		onRateLimit = func(r RateLimitInfo) {
			s.client.rateLimit = r
			s.client.onRateLimit(r)
		}
	} else {
		onRateLimit = func(r RateLimitInfo) { s.client.rateLimit = r }
	}
	out, _, err := s.client.postWithRetry(ctx, path, body, onRateLimit)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Email        string `json:"email"`
		Domain       string `json:"domain"`
		IsDisposable bool   `json:"is_disposable"`
	}
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, err
	}
	return &EmailCheckResult{
		Email:        raw.Email,
		Domain:       raw.Domain,
		IsDisposable: raw.IsDisposable,
	}, nil
}

// IsDisposable returns true if the email domain is disposable.
func (s *EmailService) IsDisposable(ctx context.Context, email string) (bool, error) {
	r, err := s.Check(ctx, email)
	if err != nil {
		return false, err
	}
	return r.IsDisposable, nil
}
