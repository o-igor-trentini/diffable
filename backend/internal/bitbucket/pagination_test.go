package bitbucket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchAllPages_MultiPage(t *testing.T) {
	var requestCount atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := requestCount.Add(1)

		var resp PaginatedResponse[Commit]
		switch page {
		case 1:
			resp = PaginatedResponse[Commit]{
				Values:  []Commit{{Hash: "aaa"}, {Hash: "bbb"}},
				PageLen: 2,
			}
			resp.Next = "http://" + r.Host + "/page2"
		case 2:
			resp = PaginatedResponse[Commit]{
				Values:  []Commit{{Hash: "ccc"}, {Hash: "ddd"}},
				PageLen: 2,
			}
			resp.Next = "http://" + r.Host + "/page3"
		case 3:
			resp = PaginatedResponse[Commit]{
				Values:  []Commit{{Hash: "eee"}},
				PageLen: 1,
			}
		}

		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	commits, err := fetchAllPages[Commit](context.Background(), httpClient, srv.URL+"/page1", "Basic dGVzdA==")

	require.NoError(t, err)
	assert.Len(t, commits, 5)
	assert.Equal(t, "aaa", commits[0].Hash)
	assert.Equal(t, "eee", commits[4].Hash)
	assert.Equal(t, int32(3), requestCount.Load())
}

func TestFetchAllPages_EmptyFirstPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(PaginatedResponse[Commit]{
			Values:  []Commit{},
			PageLen: 0,
		})
	}))
	defer srv.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	commits, err := fetchAllPages[Commit](context.Background(), httpClient, srv.URL, "Basic dGVzdA==")

	require.NoError(t, err)
	assert.Empty(t, commits)
}

func TestFetchAllPages_NoNextOnFirstPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(PaginatedResponse[Commit]{
			Values:  []Commit{{Hash: "only"}},
			PageLen: 1,
			Size:    1,
		})
	}))
	defer srv.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	commits, err := fetchAllPages[Commit](context.Background(), httpClient, srv.URL, "Basic dGVzdA==")

	require.NoError(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "only", commits[0].Hash)
}

func TestFetchAllPages_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := PaginatedResponse[Commit]{
			Values: []Commit{{Hash: "aaa"}},
		}
		resp.Next = "http://" + r.Host + "/page2"
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	_, err := fetchAllPages[Commit](ctx, httpClient, srv.URL, "Basic dGVzdA==")

	require.Error(t, err)
}

func TestFetchAllPages_ErrorOnSecondPage(t *testing.T) {
	var requestCount atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := requestCount.Add(1)
		if page == 1 {
			resp := PaginatedResponse[Commit]{
				Values: []Commit{{Hash: "aaa"}},
			}
			resp.Next = "http://" + r.Host + "/page2"
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer srv.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	_, err := fetchAllPages[Commit](context.Background(), httpClient, srv.URL, "Basic dGVzdA==")

	require.Error(t, err)
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
}
