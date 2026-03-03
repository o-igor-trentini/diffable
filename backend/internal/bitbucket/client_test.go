package bitbucket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := NewClient(Config{
		BaseURL:  srv.URL,
		Email:    "user@example.com",
		APIToken: "test-token",
		Timeout:  5 * time.Second,
	})
	return srv, client
}

func expectedAuthHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte("user@example.com:test-token"))
}

func TestClient_AuthHeader(t *testing.T) {
	var gotAuth string
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(Commit{Hash: "abc123"})
	})

	_, err := client.GetCommit(context.Background(), "ws", "repo", "abc123")
	require.NoError(t, err)
	assert.Equal(t, expectedAuthHeader(), gotAuth)
}

func TestClient_GetCommit(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/myws/myrepo/commit/abc123", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		json.NewEncoder(w).Encode(Commit{
			Hash:    "abc123",
			Message: "fix: something",
			Date:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Author:  CommitAuthor{Raw: "User <user@example.com>"},
		})
	})

	commit, err := client.GetCommit(context.Background(), "myws", "myrepo", "abc123")
	require.NoError(t, err)
	assert.Equal(t, "abc123", commit.Hash)
	assert.Equal(t, "fix: something", commit.Message)
	assert.Equal(t, "User <user@example.com>", commit.Author.Raw)
}

func TestClient_GetCommitDiff(t *testing.T) {
	diffContent := "diff --git a/file.go b/file.go\n--- a/file.go\n+++ b/file.go\n@@ -1,3 +1,4 @@\n+new line"

	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/ws/repo/diff/abc123", r.URL.Path)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, diffContent)
	})

	diff, err := client.GetCommitDiff(context.Background(), "ws", "repo", "abc123")
	require.NoError(t, err)
	assert.Equal(t, diffContent, diff)
}

func TestClient_GetCommitDiff_LargeDiff(t *testing.T) {
	largeDiff := make([]byte, 1024*1024)
	for i := range largeDiff {
		largeDiff[i] = 'x'
	}

	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(largeDiff)
	})

	diff, err := client.GetCommitDiff(context.Background(), "ws", "repo", "abc123")
	require.NoError(t, err)
	assert.Len(t, diff, 1024*1024)
}

func TestClient_GetDiffstat(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/ws/repo/diffstat/abc123", r.URL.Path)
		json.NewEncoder(w).Encode(PaginatedResponse[DiffstatEntry]{
			Values: []DiffstatEntry{
				{Status: "modified", LinesAdded: 10, LinesRemoved: 5},
				{Status: "added", LinesAdded: 20, LinesRemoved: 0},
			},
			PageLen: 2,
			Size:    2,
		})
	})

	resp, err := client.GetDiffstat(context.Background(), "ws", "repo", "abc123")
	require.NoError(t, err)
	assert.Len(t, resp.Values, 2)
	assert.Equal(t, "modified", resp.Values[0].Status)
	assert.Equal(t, 10, resp.Values[0].LinesAdded)
	assert.Equal(t, 5, resp.Values[0].LinesRemoved)
}

func TestClient_ListCommitsInRange(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/ws/repo/commits", r.URL.Path)
		assert.Equal(t, "hash-include", r.URL.Query().Get("include"))
		assert.Equal(t, "hash-exclude", r.URL.Query().Get("exclude"))
		json.NewEncoder(w).Encode(PaginatedResponse[Commit]{
			Values: []Commit{
				{Hash: "aaa"},
				{Hash: "bbb"},
			},
		})
	})

	commits, err := client.ListCommitsInRange(context.Background(), "ws", "repo", "hash-include", "hash-exclude")
	require.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "aaa", commits[0].Hash)
	assert.Equal(t, "bbb", commits[1].Hash)
}

func TestClient_GetPullRequest(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/ws/repo/pullrequests/42", r.URL.Path)
		json.NewEncoder(w).Encode(PullRequest{
			ID:    42,
			Title: "My PR",
			State: "OPEN",
			Source: PRRef{Branch: Branch{Name: "feature"}},
			Destination: PRRef{Branch: Branch{Name: "main"}},
		})
	})

	pr, err := client.GetPullRequest(context.Background(), "ws", "repo", 42)
	require.NoError(t, err)
	assert.Equal(t, 42, pr.ID)
	assert.Equal(t, "My PR", pr.Title)
	assert.Equal(t, "OPEN", pr.State)
	assert.Equal(t, "feature", pr.Source.Branch.Name)
	assert.Equal(t, "main", pr.Destination.Branch.Name)
}

func TestClient_GetPullRequestDiff(t *testing.T) {
	diffContent := "diff --git a/main.go b/main.go\n+added line"

	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/ws/repo/pullrequests/10/diff", r.URL.Path)
		fmt.Fprint(w, diffContent)
	})

	diff, err := client.GetPullRequestDiff(context.Background(), "ws", "repo", 10)
	require.NoError(t, err)
	assert.Equal(t, diffContent, diff)
}

func TestClient_GetPullRequestCommits(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/ws/repo/pullrequests/10/commits", r.URL.Path)
		json.NewEncoder(w).Encode(PaginatedResponse[Commit]{
			Values: []Commit{{Hash: "c1"}, {Hash: "c2"}},
		})
	})

	commits, err := client.GetPullRequestCommits(context.Background(), "ws", "repo", 10)
	require.NoError(t, err)
	assert.Len(t, commits, 2)
}

func TestClient_ListRepositories(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repositories/myws", r.URL.Path)
		json.NewEncoder(w).Encode(PaginatedResponse[Repository]{
			Values: []Repository{
				{Slug: "repo1", Name: "Repo One", FullName: "myws/repo1"},
				{Slug: "repo2", Name: "Repo Two", FullName: "myws/repo2"},
			},
		})
	})

	repos, err := client.ListRepositories(context.Background(), "myws")
	require.NoError(t, err)
	assert.Len(t, repos, 2)
	assert.Equal(t, "repo1", repos[0].Slug)
}

func TestClient_Error401(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "invalid credentials")
	})

	_, err := client.GetCommit(context.Background(), "ws", "repo", "abc")
	require.Error(t, err)
	var unauthErr *UnauthorizedError
	assert.ErrorAs(t, err, &unauthErr)
}

func TestClient_Error404(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "not found")
	})

	_, err := client.GetCommit(context.Background(), "ws", "repo", "nonexistent")
	require.Error(t, err)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestClient_Error429(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "120")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	_, err := client.GetCommit(context.Background(), "ws", "repo", "abc")
	require.Error(t, err)
	var rlErr *RateLimitedError
	assert.ErrorAs(t, err, &rlErr)
	assert.Equal(t, 120*time.Second, rlErr.RetryAfter)
}

func TestClient_Error500(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal error")
	})

	_, err := client.GetCommit(context.Background(), "ws", "repo", "abc")
	require.Error(t, err)
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 500, apiErr.StatusCode)
}

func TestClient_ContextCancelled(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		json.NewEncoder(w).Encode(Commit{Hash: "abc"})
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetCommit(ctx, "ws", "repo", "abc")
	require.Error(t, err)
}
