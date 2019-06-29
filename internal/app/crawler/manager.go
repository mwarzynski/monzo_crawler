package crawler

import (
	"context"
	ohttp "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mwarzynski/crawler/internal/app/crawler/http"
	"github.com/mwarzynski/crawler/internal/app/crawler/queue"
	"github.com/mwarzynski/crawler/internal/app/crawler/sitemap"
	"github.com/mwarzynski/crawler/pkg/logging"
)

type Manager struct {
	processorWorkers int

	queue            queue.FIFO
	history          *history
	sitemapGenerator *sitemap.Generator
	fetcherCreator   http.FetcherCreator

	baseURL               url.URL
	disallowedURLPrefixes []string
	log                   logging.Logger
}

func NewManager(processorWorkers int, baseURL url.URL, fetcherCreator http.FetcherCreator, log logging.Logger) *Manager {
	return &Manager{
		processorWorkers: processorWorkers,
		queue:            queue.NewFIFOSlice(100),
		history:          newHistory(),
		sitemapGenerator: sitemap.NewGenerator(),
		fetcherCreator:   fetcherCreator,
		baseURL:          baseURL,
		log:              logging.WithFields(log, "crawler", "manager"),
	}
}

func (m *Manager) SitemapGenerator(ctx context.Context) (*sitemap.Generator, error) {
	disallowPrefixes, err := m.fetchRobotsRules(ctx)
	if err != nil {
		m.log.Errorf("fetching robots rules: %s", err)
	}
	m.disallowedURLPrefixes = disallowPrefixes

	workers := m.processorWorkers
	availableWorkers := workers

	jobs := make(chan job)
	jobResults := make(chan jobResult)

	gCtx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel() // exit processors when we no longer need them
	m.initializeProcessors(gCtx, jobs, jobResults, workers)

	m.addURL(m.baseURL)

	for {
		// Try to pop the URL to process.
		u, urlToProcessExists := m.queue.Pop()
		// If there are no URLs to process and all workers are idle, then this is the end.
		if !urlToProcessExists && workers == availableWorkers {
			break
		}

		// Send / Receive the Job.
		var result *jobResult
		workersChange := 0
		if !urlToProcessExists {
			result, workersChange = m.waitForResult(gCtx, jobResults)
		} else {
			result, workersChange = m.scheduleJobOrWaitForResult(gCtx, job{url: u}, jobs, jobResults)
		}
		if result != nil {
			m.handleResult(*result)
		}

		// Update the workers count or exit.
		select {
		case <-gCtx.Done():
			return nil, context.Canceled
		default:
			availableWorkers += workersChange
		}
	}

	return m.sitemapGenerator, nil
}

func (m *Manager) initializeProcessors(
	ctx context.Context,
	jobs <-chan job,
	jobResults chan<- jobResult,
	workers int,
) {
	fetcher := m.fetcherCreator()
	for i := 0; i < workers; i++ {
		processor := newProcessor(fetcher, m.baseURL, m.log)
		_ = processor.Run(ctx, jobs, jobResults)
	}
}

func (m *Manager) scheduleJobOrWaitForResult(
	ctx context.Context,
	j job,
	jobs chan<- job,
	jobResults <-chan jobResult,
) (result *jobResult, workerChange int) {
	select {
	case jobs <- j:
		return nil, -1
	case r := <-jobResults:
		return &r, 1
	case <-ctx.Done():
		return nil, 0
	}
}

func (m *Manager) waitForResult(ctx context.Context, jobResults <-chan jobResult) (*jobResult, int) {
	select {
	case result := <-jobResults:
		return &result, 1
	case <-ctx.Done():
		return nil, 0
	}
}

func (m *Manager) handleResult(result jobResult) {
	if result.statusCode == ohttp.StatusNotFound {
		return
	}
	if result.statusCode != ohttp.StatusOK || result.err != nil {
		m.log.Infof("fetching, status code=%d, err: %s", result.statusCode, result.err)
	}
	for _, u := range result.urls {
		m.addURL(u)
	}
}

func (m *Manager) addURL(url url.URL) {
	if url.Host != m.baseURL.Host {
		return
	}
	urlRaw := url.String()
	if !strings.HasPrefix(urlRaw, m.baseURL.String()) {
		return
	}
	for _, prefix := range m.disallowedURLPrefixes {
		if strings.HasPrefix(urlRaw, prefix) {
			return
		}
	}
	if m.history.URLWasAlreadyProcessed(url) {
		return
	}
	m.sitemapGenerator.AddEntry(sitemap.Entry{
		Location: url,
	})
	m.history.SetURLProcessed(url)
	m.queue.Push(url)
}

func (m *Manager) fetchRobotsRules(ctx context.Context) ([]string, error) {
	robotsURL := m.baseURL
	robotsURL.Path += "/robots.txt"
	body, statusCode, err := m.fetcherCreator().Fetch(ctx, robotsURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Status Code=%d", statusCode)
	}

	disallows := make([]string, 0)
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "Disallow: ") {
			continue
		}
		disallowURL := m.baseURL
		disallowURL.Path = strings.Replace(line, "Disallow: ", "", 1)
		disallows = append(disallows, disallowURL.String())
	}

	return disallows, nil
}
