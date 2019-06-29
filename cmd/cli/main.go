package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/mwarzynski/crawler/internal/adapter/fetcher"
	"github.com/mwarzynski/crawler/internal/app"
	"github.com/mwarzynski/crawler/internal/app/crawler/http"
	"github.com/mwarzynski/crawler/internal/app/crawler/sitemap"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	log.Info("Hello, I am your crawler!")

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("You need to pass the URL as a first argument.")
	}

	urlRaw := args[0]
	var sitemapType sitemap.Type = sitemap.TypePlaintext
	if len(args) > 1 {
		switch sitemap.Type(args[1]) {
		case sitemap.TypePlaintext:
			sitemapType = sitemap.TypePlaintext
		case sitemap.TypeXML:
			sitemapType = sitemap.TypeXML
		default:
			log.Fatalf("Invalid sitemap type.")
		}
	}

	fetcherCreator := func() http.Fetcher {
		return fetcher.NewHTTPClient("crawler-bot", time.Minute, log)
	}
	service := app.NewService(fetcherCreator, log)

	u, err := url.Parse(urlRaw)
	if err != nil {
		log.Fatalf("couldn't parse url '%s' err: %s", urlRaw, err)
	}
	sitemap, err := service.GenerateSitemap(context.Background(), *u, sitemapType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", string(sitemap))
}
