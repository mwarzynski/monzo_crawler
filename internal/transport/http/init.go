package http

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/mwarzynski/crawler/internal/app"
	"github.com/mwarzynski/crawler/pkg/logging"
)

func Init(addr string, service *app.Service, log logging.Logger) error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/sitemap", HandleSitemap(service, log))

	// TODO: We would set up 'management' server for these endpoints (that would serve on a different port).
	// In case prometheus is used for metrics, we would add /metrics at this secondary server as well.
	r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	r.Get("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Errorf("pprof ListenAndServe: %s", err.Error())
		}
	}()
	return http.ListenAndServe(addr, r)
}
