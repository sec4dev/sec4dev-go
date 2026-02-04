package sec4dev

import (
	"testing"
)

func TestNewClient_ValidAPIKey(t *testing.T) {
	client, err := NewClient("sec4_test_key")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
	if client.APIKey != "sec4_test_key" {
		t.Errorf("APIKey = %q, want sec4_test_key", client.APIKey)
	}
}

func TestNewClient_RejectsEmptyKey(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestNewClient_RejectsInvalidPrefix(t *testing.T) {
	_, err := NewClient("invalid_key")
	if err == nil {
		t.Fatal("expected error for key without sec4_")
	}
}

func TestNewClient_WithBaseURL(t *testing.T) {
	client, err := NewClient("sec4_k", WithBaseURL("https://custom.example.com/v1"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.BaseURL != "https://custom.example.com/v1" {
		t.Errorf("BaseURL = %q", client.BaseURL)
	}
}

