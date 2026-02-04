# Sec4Dev Go SDK

Go client for the [Sec4Dev Security Checks API](https://api.sec4.dev): disposable email detection and IP classification.

## Install

```bash
go get github.com/sec4dev/sec4dev-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sec4dev/sec4dev-go"
)

func main() {
	client, err := sec4dev.NewClient("sec4_your_api_key")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// Email check
	result, err := client.Email().Check(ctx, "user@tempmail.com")
	if err != nil {
		if _, ok := err.(*sec4dev.ValidationError); ok {
			log.Printf("Invalid email: %v", err)
			return
		}
		log.Fatal(err)
	}
	if result.IsDisposable {
		fmt.Printf("Blocked: %s is a disposable domain\n", result.Domain)
	}

	// IP check
	ipResult, err := client.IP().Check(ctx, "203.0.113.42")
	if err != nil {
		if rl, ok := err.(*sec4dev.RateLimitError); ok {
			log.Printf("Rate limited. Retry in %ds", rl.RetryAfter)
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("IP Type: %s\n", ipResult.Classification)
	fmt.Printf("Confidence: %.0f%%\n", ipResult.Confidence*100)
	if ipResult.Signals.IsHosting {
		fmt.Printf("Hosting provider: %s\n", ipResult.Network.Provider)
	}
}
```

## Options

Use functional options when creating the client:

- `sec4dev.WithBaseURL(url)` — API base URL (default: `https://api.sec4.dev/api/v1`)
- `sec4dev.WithRetries(n)` — Retry attempts (default: 3)
- `sec4dev.WithRetryDelay(ms)` — Base retry delay in ms (default: 1000)
- `sec4dev.WithHTTPClient(hc)` — Custom `*http.Client`
- `sec4dev.WithRateLimitCallback(fn)` — Callback for rate limit updates