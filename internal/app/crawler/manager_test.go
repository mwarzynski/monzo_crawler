package crawler

import (
	"context"
	"fmt"
	ohttp "net/http"
	"net/url"
	"testing"

	"github.com/mwarzynski/crawler/internal/app/crawler/http"
	"github.com/sirupsen/logrus"
)

type mockFetcher struct {
	baseURL            url.URL
	disallowedPrefixes []string
	urls               map[string][]string
}

func (mf *mockFetcher) Fetch(ctx context.Context, url url.URL) ([]byte, int, error) {
	if mf.baseURL.String()+"/robots.txt" == url.String() {
		return mf.fetchRobots()
	}
	return mf.fetchSite(url)
}

func (mf *mockFetcher) fetchRobots() ([]byte, int, error) {
	robots := ""
	for _, disallowedPrefix := range mf.disallowedPrefixes {
		robots += fmt.Sprintf("Disallow: %s\n", disallowedPrefix)
	}
	return []byte(robots), ohttp.StatusOK, nil
}

func (mf *mockFetcher) fetchSite(url url.URL) ([]byte, int, error) {
	generateHTML := func(urls []string) []byte {
		html := "<html><body>"
		for _, url := range urls {
			html += fmt.Sprintf(`<a href="%s">%s</a>`, url, url)
		}
		html += "</body></html>"
		return []byte(html)
	}
	return generateHTML(mf.urls[url.String()]), ohttp.StatusOK, nil
}

func TestManager(t *testing.T) {
	log := logrus.New()
	tests := []struct {
		name             string              // Name of the test.
		baseURL          string              // BaseURL that user provides as input.
		pageLinks        map[string][]string // Map: url -> urls; Graph of our site.
		robotsDisallowed []string            // Robots functionality, entry: 'Disallow: prefix'
		expectedLinks    []string            // Expected output links (that goes to sitemap).
	}{
		{
			name:    "simple site with only one page",
			baseURL: "https://google.com",
			pageLinks: map[string][]string{
				"https://google.com": []string{},
			},
			expectedLinks: []string{"https://google.com"},
		},
		{
			name:    "site with two sites that point to each other",
			baseURL: "https://google.com",
			pageLinks: map[string][]string{
				"https://google.com": []string{
					"https://google.com/1",
				},
				"https://google.com/1": []string{
					"https://google.com",
				},
			},
			expectedLinks: []string{
				"https://google.com",
				"https://google.com/1",
			},
		},
		{
			name:    "site with two sites that point to each other, but the second one is disallowed by robots",
			baseURL: "https://google.com",
			robotsDisallowed: []string{
				"/1",
			},
			pageLinks: map[string][]string{
				"https://google.com": []string{
					"https://google.com/1",
				},
				"https://google.com/1": []string{
					"https://google.com",
				},
			},
			expectedLinks: []string{
				"https://google.com",
			},
		},
		// It should be relatively easy, to test new site configurations / functionalities.
		// I much more like the approach of testing the upper level, than single components.
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			baseURL, err := url.Parse(test.baseURL)
			if err != nil {
				t.Fatalf("invalid base URL '%s': %s", test.baseURL, err)
			}
			fetcherCreator := func() http.Fetcher {
				return &mockFetcher{
					disallowedPrefixes: test.robotsDisallowed,
					baseURL:            *baseURL,
					urls:               test.pageLinks,
				}
			}
			manager := NewManager(3, *baseURL, fetcherCreator, log)

			ctx := context.Background()
			sg, err := manager.SitemapGenerator(ctx)
			if err != nil {
				t.Fatalf("couldn't generate sitemap: %s", err)
			}

			if len(sg.Entries) != len(test.expectedLinks) {
				t.Fatalf("received invalid number of links: got: %d, want: %d\nsitemap: %v",
					len(sg.Entries), len(test.expectedLinks), sg.Entries)
			}
			for i := range sg.Entries {
				if sg.Entries[i].Location.String() != test.expectedLinks[i] {
					t.Errorf("received invalid link(i=%d): got: %q, want: %q\nsitemap: %v",
						i, sg.Entries[i].Location.String(), test.expectedLinks[i], sg.Entries)
				}
			}
		})
	}
}
