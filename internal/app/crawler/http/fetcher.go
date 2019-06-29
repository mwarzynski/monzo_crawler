package http

import (
	"context"
	"errors"
	"net/url"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

// Fetcher provides functionality of fetching and rendering the contents of web sites.
// In the simple approach it might be just the HTTP client. However, you could also provide here an implementation that
// uses full headless web browser for rendering sites (especially modern ones).
type Fetcher interface {
	Fetch(ctx context.Context, url url.URL) ([]byte, int, error)
}

type FetcherCreator = func() Fetcher
