package sec4dev

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmailCheck_ReturnsResult(t *testing.T) {
	body := map[string]interface{}{
		"email":         "user@tempmail.com",
		"domain":        "tempmail.com",
		"is_disposable": true,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/email/check" || r.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}))
	defer server.Close()

	client, err := NewClient("sec4_test", WithBaseURL(server.URL+"/api/v1"), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ctx := context.Background()

	result, err := client.Email().Check(ctx, "user@tempmail.com")
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if result.Email != "user@tempmail.com" || result.Domain != "tempmail.com" || !result.IsDisposable {
		t.Errorf("result = %+v", result)
	}
}

func TestEmailIsDisposable_True(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"email": "x@disposable.com", "domain": "disposable.com", "is_disposable": true,
		})
	}))
	defer server.Close()

	client, _ := NewClient("sec4_test", WithBaseURL(server.URL+"/api/v1"), WithHTTPClient(server.Client()))
	ctx := context.Background()

	ok, err := client.Email().IsDisposable(ctx, "x@disposable.com")
	if err != nil {
		t.Fatalf("IsDisposable: %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestEmailIsDisposable_False(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"email": "user@gmail.com", "domain": "gmail.com", "is_disposable": false,
		})
	}))
	defer server.Close()

	client, _ := NewClient("sec4_test", WithBaseURL(server.URL+"/api/v1"), WithHTTPClient(server.Client()))
	ctx := context.Background()

	ok, err := client.Email().IsDisposable(ctx, "user@gmail.com")
	if err != nil {
		t.Fatalf("IsDisposable: %v", err)
	}
	if ok {
		t.Error("expected false")
	}
}

func TestEmailCheck_ValidatesInput(t *testing.T) {
	client, _ := NewClient("sec4_test")
	ctx := context.Background()

	_, err := client.Email().Check(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
	_, err = client.Email().Check(ctx, "not-an-email")
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
}
