package sec4dev

import (
	"context"
	"encoding/json"
	"strings"
)

// IPService provides IP check operations.
type IPService struct {
	client *Client
}

// Check classifies an IP address.
func (s *IPService) Check(ctx context.Context, ip string) (*IPCheckResult, error) {
	if err := ValidateIP(ip); err != nil {
		return nil, err
	}
	path := "/ip/check"
	body := map[string]string{"ip": strings.TrimSpace(ip)}
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
		IP             string  `json:"ip"`
		Classification string  `json:"classification"`
		Confidence     float64 `json:"confidence"`
		Signals        struct {
			IsHosting     bool `json:"is_hosting"`
			IsResidential  bool `json:"is_residential"`
			IsMobile      bool `json:"is_mobile"`
			IsVPN         bool `json:"is_vpn"`
			IsTor         bool `json:"is_tor"`
			IsProxy       bool `json:"is_proxy"`
		} `json:"signals"`
		Network struct {
			ASN      *int   `json:"asn"`
			Org      string `json:"org"`
			Provider string `json:"provider"`
		} `json:"network"`
		Geo struct {
			Country string `json:"country"`
			Region  string `json:"region"`
		} `json:"geo"`
	}
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, err
	}
	return &IPCheckResult{
		IP:             raw.IP,
		Classification: raw.Classification,
		Confidence:     raw.Confidence,
		Signals: IPSignals{
			IsHosting:     raw.Signals.IsHosting,
			IsResidential: raw.Signals.IsResidential,
			IsMobile:      raw.Signals.IsMobile,
			IsVPN:         raw.Signals.IsVPN,
			IsTor:         raw.Signals.IsTor,
			IsProxy:       raw.Signals.IsProxy,
		},
		Network: IPNetwork{
			ASN:      raw.Network.ASN,
			Org:      raw.Network.Org,
			Provider: raw.Network.Provider,
		},
		Geo: IPGeo{
			Country: raw.Geo.Country,
			Region:  raw.Geo.Region,
		},
	}, nil
}

// IsHosting returns true if the IP is classified as hosting.
func (s *IPService) IsHosting(ctx context.Context, ip string) (bool, error) {
	r, err := s.Check(ctx, ip)
	if err != nil {
		return false, err
	}
	return r.Signals.IsHosting, nil
}

// IsVPN returns true if the IP is classified as VPN.
func (s *IPService) IsVPN(ctx context.Context, ip string) (bool, error) {
	r, err := s.Check(ctx, ip)
	if err != nil {
		return false, err
	}
	return r.Signals.IsVPN, nil
}

// IsTor returns true if the IP is classified as TOR.
func (s *IPService) IsTor(ctx context.Context, ip string) (bool, error) {
	r, err := s.Check(ctx, ip)
	if err != nil {
		return false, err
	}
	return r.Signals.IsTor, nil
}

// IsResidential returns true if the IP is classified as residential.
func (s *IPService) IsResidential(ctx context.Context, ip string) (bool, error) {
	r, err := s.Check(ctx, ip)
	if err != nil {
		return false, err
	}
	return r.Signals.IsResidential, nil
}

// IsMobile returns true if the IP is classified as mobile.
func (s *IPService) IsMobile(ctx context.Context, ip string) (bool, error) {
	r, err := s.Check(ctx, ip)
	if err != nil {
		return false, err
	}
	return r.Signals.IsMobile, nil
}
