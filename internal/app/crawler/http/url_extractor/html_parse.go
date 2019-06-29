package url_extractor

import (
	"bytes"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type HTMLParse struct{}

func NewHTMLParse() *HTMLParse {
	return &HTMLParse{}
}

func (uer *HTMLParse) ExtractURLs(baseURL url.URL, body []byte) ([]url.URL, error) {
	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "parsing body")
	}

	urls := uer.extractURLsFromAHrefs(root)
	resolvedURLs := uer.resolveURLs(baseURL, urls)
	urls = uer.normalizeURLs(resolvedURLs)

	return urls, nil
}

func (uer *HTMLParse) extractURLsFromAHrefs(root *html.Node) []url.URL {
	urls := make([]url.URL, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}
					urls = append(urls, *u)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
	return urls
}

func (uer *HTMLParse) resolveURLs(base url.URL, links []url.URL) []url.URL {
	newLinks := make([]url.URL, 0, len(links))
	for _, link := range links {
		linkResolved := base.ResolveReference(&link)
		newLinks = append(newLinks, *linkResolved)
	}
	return newLinks
}

func (uer *HTMLParse) normalizeURLs(urls []url.URL) []url.URL {
	normalizedURLs := make([]url.URL, 0, len(urls))
	for _, u := range urls {
		normalizedURLs = append(normalizedURLs, normalizeURL(u))
	}
	return normalizedURLs
}
