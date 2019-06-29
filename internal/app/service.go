package app

import (
	"context"
	"net/url"

	"github.com/mwarzynski/crawler/pkg/logging"

	"github.com/mwarzynski/crawler/internal/app/crawler"
	"github.com/mwarzynski/crawler/internal/app/crawler/http"
	"github.com/mwarzynski/crawler/internal/app/crawler/sitemap"
)

// processorsCount should be configurable. Maybe some sites need more workers than others.
// For now, let's leave it as a const.
const processorsCount = 10

type Service struct {
	fetcherCreator http.FetcherCreator

	log logging.Logger
}

func NewService(fetcherCreator http.FetcherCreator, log logging.Logger) *Service {
	return &Service{
		fetcherCreator: fetcherCreator,
		log:            logging.WithFields(log, "app", "service"),
	}
}

func (s *Service) GenerateSitemap(ctx context.Context, baseURL url.URL, sitemapType sitemap.Type) ([]byte, error) {
	manager := crawler.NewManager(
		processorsCount,
		baseURL,
		s.fetcherCreator,
		s.log.WithField("url", baseURL.String()),
	)
	sitemapGenerator, err := manager.SitemapGenerator(ctx)
	if err != nil {
		return []byte{}, err
	}
	return sitemapGenerator.Generate(sitemapType)
}
