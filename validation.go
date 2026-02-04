package sec4dev

import (
	"net"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// ValidateEmail returns a ValidationError if the email is invalid.
func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{baseError("Email is required", 422, nil)}
	}
	s := strings.TrimSpace(email)
	if s == "" {
		return &ValidationError{baseError("Email cannot be empty", 422, nil)}
	}
	if !emailRegex.MatchString(s) {
		return &ValidationError{baseError("Invalid email format", 422, nil)}
	}
	return nil
}

// ValidateIP returns a ValidationError if the IP is invalid.
func ValidateIP(ip string) error {
	if ip == "" {
		return &ValidationError{baseError("IP address is required", 422, nil)}
	}
	s := strings.TrimSpace(ip)
	if s == "" {
		return &ValidationError{baseError("IP address cannot be empty", 422, nil)}
	}
	if net.ParseIP(s) == nil {
		return &ValidationError{baseError("Invalid IP address format", 422, nil)}
	}
	return nil
}
