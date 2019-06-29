package crawler

import (
	"context"
	"net/url"
	"sync"
	"testing"

	"github.com/mwarzynski/crawler/internal/app/crawler/http"
	"github.com/sirupsen/logrus"
)

func TestProcessorOnContextCancel(t *testing.T) {
	// Processor's context must be supplied with the timeout.
	// We do want to check if the processors are going to exit in this case (as not to leak the resources -- goroutines
	// also have the upper limit).

	log := logrus.New()
	fetcherCreator := func() http.Fetcher {
		return nil
	}
	u, _ := url.Parse("https://google.com")
	processor := newProcessor(fetcherCreator(), *u, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobsChan := make(chan job)
	jobsResults := make(chan jobResult)

	processorFinished := false
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-processor.Run(ctx, jobsChan, jobsResults)
		processorFinished = true
		wg.Done()
	}()

	cancel()
	wg.Wait()

	if !processorFinished {
		t.Errorf("processor didn't finish after context cancel")
	}
}
