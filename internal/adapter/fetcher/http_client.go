package fetcher

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	chttp "github.com/mwarzynski/crawler/internal/app/crawler/http"
	"github.com/mwarzynski/crawler/pkg/logging"
	"github.com/pkg/errors"
)

type HTTPClient struct {
	httpDoer *http.Client

	name string
	log  logging.Logger
}

func NewHTTPClient(name string, timeout time.Duration, log logging.Logger) *HTTPClient {
	return &HTTPClient{
		httpDoer: &http.Client{
			Timeout: timeout,
		},
		name: name,
		log:  logging.WithFields(log, "fetcher", "HTTPClient"),
	}
}

func (s *HTTPClient) Fetch(ctx context.Context, url url.URL) ([]byte, int, error) {
	s.log.Debugf("Fetching URL=%s", url.String())
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "couldn't create request")
	}
	req.Header.Set("User-Agent", s.name)
	resp, err := s.httpDoer.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "couldn't do HTTP request")
	}
	if resp == nil {
		return nil, 0, errors.Errorf("body is nil")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, chttp.ErrInvalidStatusCode
	}
	if resp.Body == nil {
		return nil, resp.StatusCode, errors.Errorf("body is nil")
	}
	reader := io.LimitReader(resp.Body, 10*1024*1024) // Limit reading body to 10MB.
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "reading request body")
	}
	return data, resp.StatusCode, nil
}
