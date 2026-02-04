package sec4dev

// EmailCheckResult is the result of an email check.
type EmailCheckResult struct {
	Email        string `json:"email"`
	Domain       string `json:"domain"`
	IsDisposable bool   `json:"is_disposable"`
}

// IPSignals holds signals from an IP check.
type IPSignals struct {
	IsHosting     bool `json:"is_hosting"`
	IsResidential bool `json:"is_residential"`
	IsMobile      bool `json:"is_mobile"`
	IsVPN         bool `json:"is_vpn"`
	IsTor         bool `json:"is_tor"`
	IsProxy       bool `json:"is_proxy"`
}

// IPNetwork holds network info from an IP check.
type IPNetwork struct {
	ASN      *int   `json:"asn"`
	Org      string `json:"org,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// IPGeo holds geo info from an IP check.
type IPGeo struct {
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
}

// IPCheckResult is the result of an IP check.
type IPCheckResult struct {
	IP             string    `json:"ip"`
	Classification string    `json:"classification"`
	Confidence     float64   `json:"confidence"`
	Signals        IPSignals `json:"signals"`
	Network        IPNetwork `json:"network"`
	Geo            IPGeo     `json:"geo"`
}

// RateLimitInfo holds rate limit data from response headers.
type RateLimitInfo struct {
	Limit        int
	Remaining    int
	ResetSeconds int
}
