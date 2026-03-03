package bitbucket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client interface {
	GetCommit(ctx context.Context, workspace, repoSlug, hash string) (*Commit, error)
	GetCommitDiff(ctx context.Context, workspace, repoSlug, spec string) (string, error)
	GetDiffstat(ctx context.Context, workspace, repoSlug, spec string) (*PaginatedResponse[DiffstatEntry], error)
	ListCommitsInRange(ctx context.Context, workspace, repoSlug, include, exclude string) ([]Commit, error)
	GetPullRequest(ctx context.Context, workspace, repoSlug string, prID int) (*PullRequest, error)
	GetPullRequestDiff(ctx context.Context, workspace, repoSlug string, prID int) (string, error)
	GetPullRequestCommits(ctx context.Context, workspace, repoSlug string, prID int) ([]Commit, error)
	ListRepositories(ctx context.Context, workspace string) ([]Repository, error)
}

type Config struct {
	BaseURL  string
	Email    string
	APIToken string
	Timeout  time.Duration
}

type bitbucketClient struct {
	httpClient *http.Client
	baseURL    string
	authHeader string
}

func NewClient(cfg Config) Client {
	credentials := base64.StdEncoding.EncodeToString([]byte(cfg.Email + ":" + cfg.APIToken))
	return &bitbucketClient{
		httpClient: &http.Client{Timeout: cfg.Timeout},
		baseURL:    cfg.BaseURL,
		authHeader: "Basic " + credentials,
	}
}

func (c *bitbucketClient) GetCommit(ctx context.Context, workspace, repoSlug, hash string) (*Commit, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/commit/%s", c.baseURL, workspace, repoSlug, hash)
	var commit Commit
	if err := c.doJSON(ctx, url, &commit); err != nil {
		return nil, err
	}
	return &commit, nil
}

func (c *bitbucketClient) GetCommitDiff(ctx context.Context, workspace, repoSlug, spec string) (string, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/diff/%s", c.baseURL, workspace, repoSlug, spec)
	return c.doText(ctx, url)
}

func (c *bitbucketClient) GetDiffstat(ctx context.Context, workspace, repoSlug, spec string) (*PaginatedResponse[DiffstatEntry], error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/diffstat/%s", c.baseURL, workspace, repoSlug, spec)
	var resp PaginatedResponse[DiffstatEntry]
	if err := c.doJSON(ctx, url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *bitbucketClient) ListCommitsInRange(ctx context.Context, workspace, repoSlug, include, exclude string) ([]Commit, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/commits?include=%s&exclude=%s&pagelen=100",
		c.baseURL, workspace, repoSlug, include, exclude)
	return fetchAllPages[Commit](ctx, c.httpClient, url, c.authHeader)
}

func (c *bitbucketClient) GetPullRequest(ctx context.Context, workspace, repoSlug string, prID int) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d", c.baseURL, workspace, repoSlug, prID)
	var pr PullRequest
	if err := c.doJSON(ctx, url, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (c *bitbucketClient) GetPullRequestDiff(ctx context.Context, workspace, repoSlug string, prID int) (string, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/diff", c.baseURL, workspace, repoSlug, prID)
	return c.doText(ctx, url)
}

func (c *bitbucketClient) GetPullRequestCommits(ctx context.Context, workspace, repoSlug string, prID int) ([]Commit, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/commits?pagelen=100",
		c.baseURL, workspace, repoSlug, prID)
	return fetchAllPages[Commit](ctx, c.httpClient, url, c.authHeader)
}

func (c *bitbucketClient) ListRepositories(ctx context.Context, workspace string) ([]Repository, error) {
	url := fmt.Sprintf("%s/repositories/%s?pagelen=100", c.baseURL, workspace)
	return fetchAllPages[Repository](ctx, c.httpClient, url, c.authHeader)
}

func (c *bitbucketClient) doJSON(ctx context.Context, url string, target any) error {
	body, statusCode, err := c.doRequest(ctx, url, "application/json")
	if err != nil {
		return err
	}
	if err := checkResponseError(statusCode, body); err != nil {
		return err
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("bitbucket: decoding response: %w", err)
	}
	return nil
}

func (c *bitbucketClient) doText(ctx context.Context, url string) (string, error) {
	body, statusCode, err := c.doRequest(ctx, url, "text/plain")
	if err != nil {
		return "", err
	}
	if err := checkResponseError(statusCode, body); err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *bitbucketClient) doRequest(ctx context.Context, url, accept string) ([]byte, int, error) {
	slog.Debug("bitbucket: request", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("bitbucket: creating request: %w", err)
	}
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", accept)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("bitbucket: executing request: %w", err)
	}
	defer resp.Body.Close()

	if err := checkRateLimit(resp); err != nil {
		return nil, resp.StatusCode, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("bitbucket: reading response: %w", err)
	}

	return body, resp.StatusCode, nil
}

func checkResponseError(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return &NotFoundError{Resource: "resource"}
	case http.StatusUnauthorized:
		return &UnauthorizedError{Message: string(body)}
	default:
		if statusCode >= 400 {
			return &APIError{StatusCode: statusCode, Message: string(body)}
		}
		return nil
	}
}
