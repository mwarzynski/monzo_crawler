package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mwarzynski/crawler/internal/adapter/fetcher"
	"github.com/mwarzynski/crawler/internal/app"
	"github.com/mwarzynski/crawler/internal/app/crawler/http"
	httpAPI "github.com/mwarzynski/crawler/internal/transport/http"
)

func main() {
	log := logrus.New()

	// Create application service.
	fetcherCreator := func() http.Fetcher {
		return fetcher.NewHTTPClient("crawler-bot", time.Minute, log)
	}
	service := app.NewService(fetcherCreator, log)

	// Create HTTP server.
	listenAddr := "localhost:8000"
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		listenAddr = addr
	}
	if err := httpAPI.Init(listenAddr, service, log.WithField("component", "http")); err != nil {
		log.Errorf("http API: %s", err.Error())
	}
}
