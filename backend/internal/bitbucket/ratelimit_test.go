package bitbucket

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRateLimitHeaders(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"X-Ratelimit-Limit":     {"1000"},
			"X-Ratelimit-Nearlimit": {"false"},
		},
	}

	info := parseRateLimitHeaders(resp)
	assert.Equal(t, 1000, info.Limit)
	assert.False(t, info.NearLimit)
}

func TestParseRateLimitHeaders_NearLimit(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"X-Ratelimit-Limit":     {"1000"},
			"X-Ratelimit-Nearlimit": {"true"},
		},
	}

	info := parseRateLimitHeaders(resp)
	assert.Equal(t, 1000, info.Limit)
	assert.True(t, info.NearLimit)
}

func TestParseRateLimitHeaders_Missing(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}

	info := parseRateLimitHeaders(resp)
	assert.Equal(t, 0, info.Limit)
	assert.False(t, info.NearLimit)
}

func TestCheckRateLimit_429(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			"Retry-After": {"90"},
		},
	}

	err := checkRateLimit(resp)
	require.Error(t, err)
	var rlErr *RateLimitedError
	assert.ErrorAs(t, err, &rlErr)
	assert.Equal(t, 90*time.Second, rlErr.RetryAfter)
}

func TestCheckRateLimit_429_DefaultRetry(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{},
	}

	err := checkRateLimit(resp)
	require.Error(t, err)
	var rlErr *RateLimitedError
	assert.ErrorAs(t, err, &rlErr)
	assert.Equal(t, 60*time.Second, rlErr.RetryAfter)
}

func TestCheckRateLimit_OK(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Ratelimit-Limit":     {"1000"},
			"X-Ratelimit-Nearlimit": {"false"},
		},
	}

	err := checkRateLimit(resp)
	assert.NoError(t, err)
}

func TestCheckRateLimit_NearLimit_Warning(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Ratelimit-Limit":     {"1000"},
			"X-Ratelimit-Nearlimit": {"true"},
		},
	}

	err := checkRateLimit(resp)
	assert.NoError(t, err)
}
