package crawler

import "net/url"

type job struct {
	url url.URL
}

type jobResult struct {
	urls       []url.URL
	statusCode int
	err        error
}
