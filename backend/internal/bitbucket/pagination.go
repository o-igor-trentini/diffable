package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetchAllPages[T any](ctx context.Context, httpClient *http.Client, initialURL string, authHeader string) ([]T, error) {
	var allValues []T
	url := initialURL

	for url != "" {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating paginated request: %w", err)
		}
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Accept", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("executing paginated request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading paginated response: %w", err)
		}

		if err := checkRateLimit(resp); err != nil {
			return nil, err
		}

		if err := checkResponseError(resp.StatusCode, body); err != nil {
			return nil, err
		}

		var page PaginatedResponse[T]
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("decoding paginated response: %w", err)
		}

		allValues = append(allValues, page.Values...)
		url = page.Next
	}

	return allValues, nil
}
