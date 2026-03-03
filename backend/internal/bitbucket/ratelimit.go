package bitbucket

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type RateLimitInfo struct {
	Limit     int
	NearLimit bool
}

func parseRateLimitHeaders(resp *http.Response) RateLimitInfo {
	info := RateLimitInfo{}

	if v := resp.Header.Get("X-RateLimit-Limit"); v != "" {
		info.Limit, _ = strconv.Atoi(v)
	}

	if v := resp.Header.Get("X-RateLimit-NearLimit"); v != "" {
		info.NearLimit = v == "true"
	}

	return info
}

func checkRateLimit(resp *http.Response) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := 60 * time.Second
		if v := resp.Header.Get("Retry-After"); v != "" {
			if secs, err := strconv.Atoi(v); err == nil {
				retryAfter = time.Duration(secs) * time.Second
			}
		}
		return &RateLimitedError{RetryAfter: retryAfter}
	}

	info := parseRateLimitHeaders(resp)
	if info.NearLimit {
		slog.Warn("bitbucket: approaching rate limit",
			"limit", info.Limit,
		)
	}

	return nil
}
