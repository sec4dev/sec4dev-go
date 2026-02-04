// Package sec4dev provides the Sec4Dev Security Checks API client.
package sec4dev

// Sec4DevError is the base type for API errors.
type Sec4DevError struct {
	Message      string
	StatusCode   int
	ResponseBody interface{}
}

func (e *Sec4DevError) Error() string {
	return e.Message
}

// AuthenticationError is returned for 401 (invalid or missing API key).
type AuthenticationError struct{ *Sec4DevError }

// PaymentRequiredError is returned for 402 (quota exceeded).
type PaymentRequiredError struct{ *Sec4DevError }

// ForbiddenError is returned for 403 (account deactivated).
type ForbiddenError struct{ *Sec4DevError }

// NotFoundError is returned for 404.
type NotFoundError struct{ *Sec4DevError }

// ValidationError is returned for 422 or client-side validation failure.
type ValidationError struct{ *Sec4DevError }

// RateLimitError is returned for 429 with rate limit info.
type RateLimitError struct {
	*Sec4DevError
	RetryAfter int
	Limit      int
	Remaining  int
}

// ServerError is returned for 5xx.
type ServerError struct{ *Sec4DevError }

func baseError(message string, statusCode int, body interface{}) *Sec4DevError {
	return &Sec4DevError{
		Message:      message,
		StatusCode:   statusCode,
		ResponseBody: body,
	}
}

func errFromStatus(statusCode int, message string, body interface{}, retryAfter, limit, remaining int) error {
	base := baseError(message, statusCode, body)
	switch statusCode {
	case 401:
		return &AuthenticationError{base}
	case 402:
		return &PaymentRequiredError{base}
	case 403:
		return &ForbiddenError{base}
	case 404:
		return &NotFoundError{base}
	case 422:
		return &ValidationError{base}
	case 429:
		return &RateLimitError{Sec4DevError: base, RetryAfter: retryAfter, Limit: limit, Remaining: remaining}
	default:
		if statusCode >= 500 {
			return &ServerError{base}
		}
		return base
	}
}
