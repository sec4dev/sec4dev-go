package sec4dev

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPCheck_ReturnsResult(t *testing.T) {
	body := map[string]interface{}{
		"ip":             "203.0.113.42",
		"classification": "hosting",
		"confidence":     0.95,
		"signals": map[string]bool{
			"is_hosting": true, "is_residential": false, "is_mobile": false,
			"is_vpn": false, "is_tor": false, "is_proxy": false,
		},
		"network": map[string]interface{}{"asn": 16509, "org": "Amazon.com, Inc.", "provider": "AWS"},
		"geo":     map[string]interface{}{"country": "US", "region": nil},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/ip/check" || r.Method != http.MethodPost {
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

	result, err := client.IP().Check(ctx, "203.0.113.42")
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if result.IP != "203.0.113.42" || result.Classification != "hosting" || result.Confidence != 0.95 {
		t.Errorf("result = %+v", result)
	}
	if !result.Signals.IsHosting || result.Signals.IsVPN {
		t.Errorf("signals = %+v", result.Signals)
	}
	if result.Network.Provider != "AWS" || result.Geo.Country != "US" {
		t.Errorf("network/geo = %+v / %+v", result.Network, result.Geo)
	}
}

func TestIPIsHosting_True(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ip": "203.0.113.42", "classification": "hosting", "confidence": 0.95,
			"signals": map[string]bool{"is_hosting": true, "is_residential": false, "is_mobile": false, "is_vpn": false, "is_tor": false, "is_proxy": false},
			"network": map[string]interface{}{"asn": nil, "org": "", "provider": ""},
			"geo":     map[string]interface{}{"country": "", "region": ""},
		})
	}))
	defer server.Close()

	client, _ := NewClient("sec4_test", WithBaseURL(server.URL+"/api/v1"), WithHTTPClient(server.Client()))
	ctx := context.Background()

	ok, err := client.IP().IsHosting(ctx, "203.0.113.42")
	if err != nil {
		t.Fatalf("IsHosting: %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestIPCheck_ValidatesInput(t *testing.T) {
	client, _ := NewClient("sec4_test")
	ctx := context.Background()

	_, err := client.IP().Check(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty IP")
	}
	_, err = client.IP().Check(ctx, "not-an-ip")
	if err == nil {
		t.Fatal("expected error for invalid IP")
	}
}

func TestIPCheck_AcceptsIPv6(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ip": "::1", "classification": "unknown", "confidence": 0.0,
			"signals": map[string]bool{"is_hosting": false, "is_residential": false, "is_mobile": false, "is_vpn": false, "is_tor": false, "is_proxy": false},
			"network": map[string]interface{}{"asn": nil, "org": "", "provider": ""},
			"geo":     map[string]interface{}{"country": "", "region": ""},
		})
	}))
	defer server.Close()

	client, _ := NewClient("sec4_test", WithBaseURL(server.URL+"/api/v1"), WithHTTPClient(server.Client()))
	ctx := context.Background()

	result, err := client.IP().Check(ctx, "::1")
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if result.IP != "::1" {
		t.Errorf("IP = %q", result.IP)
	}
}
