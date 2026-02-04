package sec4dev

import (
	"testing"
)

func TestValidateEmail_AcceptsValid(t *testing.T) {
	if err := ValidateEmail("user@example.com"); err != nil {
		t.Errorf("ValidateEmail: %v", err)
	}
	if err := ValidateEmail("a@b.co"); err != nil {
		t.Errorf("ValidateEmail: %v", err)
	}
}

func TestValidateEmail_RejectsEmpty(t *testing.T) {
	if err := ValidateEmail(""); err == nil {
		t.Error("expected error for empty")
	}
	if err := ValidateEmail("   "); err == nil {
		t.Error("expected error for whitespace")
	}
}

func TestValidateEmail_RejectsInvalid(t *testing.T) {
	for _, s := range []string{"no-at-sign", "@nodomain.com", "nobody@", "a@b"} {
		if err := ValidateEmail(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}

func TestValidateIP_AcceptsIPv4(t *testing.T) {
	for _, s := range []string{"192.168.1.1", "0.0.0.0", "255.255.255.255", "203.0.113.42"} {
		if err := ValidateIP(s); err != nil {
			t.Errorf("ValidateIP(%q): %v", s, err)
		}
	}
}

func TestValidateIP_AcceptsIPv6(t *testing.T) {
	for _, s := range []string{"::1", "2001:db8::1"} {
		if err := ValidateIP(s); err != nil {
			t.Errorf("ValidateIP(%q): %v", s, err)
		}
	}
}

func TestValidateIP_RejectsEmpty(t *testing.T) {
	if err := ValidateIP(""); err == nil {
		t.Error("expected error for empty")
	}
}

func TestValidateIP_RejectsInvalid(t *testing.T) {
	for _, s := range []string{"256.1.1.1", "not.an.ip"} {
		if err := ValidateIP(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}
