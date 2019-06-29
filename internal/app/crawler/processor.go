package crawler

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/mwarzynski/crawler/internal/app/crawler/http"
	"github.com/mwarzynski/crawler/internal/app/crawler/http/url_extractor"
	"github.com/mwarzynski/crawler/pkg/logging"
)

type processor struct {
	fetcher      http.Fetcher
	urlExtractor *url_extractor.HTMLParse
	baseURL      url.URL
	log          logging.Logger
}

func newProcessor(
	fetcher http.Fetcher,
	baseURL url.URL,
	log logging.Logger,
) *processor {
	return &processor{
		fetcher:      fetcher,
		urlExtractor: url_extractor.NewHTMLParse(),
		baseURL:      baseURL,
		log:          logging.WithFields(log, "crawler", "processor"),
	}
}

func (p *processor) processJob(ctx context.Context, j job) jobResult {
	body, statusCode, err := p.fetcher.Fetch(ctx, j.url)
	if err != nil {
		return jobResult{
			statusCode: statusCode,
			err:        errors.Wrapf(err, "couldn't fetch '%s'", j.url.String()),
		}
	}
	urls, err := p.urlExtractor.ExtractURLs(j.url, body)
	if err != nil {
		return jobResult{
			statusCode: statusCode,
			err:        errors.Wrap(err, "couldn't extract urls from body"),
		}
	}
	return jobResult{
		urls:       urls,
		statusCode: statusCode,
		err:        nil,
	}
}

func (p *processor) Run(ctx context.Context, jobs <-chan job, jobResults chan<- jobResult) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		for {
			select {
			case job := <-jobs:
				result := p.processJob(ctx, job)
				jobResults <- result
			case <-ctx.Done():
				p.log.Debug("processor exited (ctx.Done)")
				close(done)
				return
			}
		}
	}()
	return done
}
