package crawler

import (
	"net/url"
)

type history struct {
	processedURLs map[string]struct{}
}

func newHistory() *history {
	return &history{
		processedURLs: make(map[string]struct{}),
	}
}

func (h *history) URLWasAlreadyProcessed(u url.URL) bool {
	_, alreadyProcessed := h.processedURLs[u.String()]
	return alreadyProcessed
}

func (h *history) SetURLProcessed(u url.URL) {
	h.processedURLs[u.String()] = struct{}{}
}
