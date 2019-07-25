package http

import (
	"net/http"
	"net/url"

	"github.com/mwarzynski/crawler/internal/app"
	"github.com/mwarzynski/crawler/internal/app/crawler/sitemap"
	"github.com/mwarzynski/crawler/pkg/logging"
)

func HandleSitemap(service *app.Service, log logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		baseURLRaw := r.URL.Query().Get("url")
		if baseURLRaw == "" {
			http.Error(w, "provided url is empty", http.StatusBadRequest)
			return
		}
		baseURL, err := url.Parse(baseURLRaw)
		if err != nil {
			http.Error(w, "provided url is invalid", http.StatusUnprocessableEntity)
			return
		}
		sitemapType := sitemap.TypePlaintext

		data, err := service.GenerateSitemap(ctx, *baseURL, sitemapType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(data); err != nil {
			log.Errorf("couldn't write data: %s", err)
		}
	}
}
