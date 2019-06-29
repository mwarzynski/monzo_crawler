package url_extractor

import (
	"net/url"
	"testing"
)

const siteWithValidHTML = `<html>
<body>
	<a href="test"/>
	<a href="/test2"/>
	<p>This is awesome</p>
	<a href="https://google.com/test"/>
</body>
</html>
`

func TestURLExtractor(t *testing.T) {
	urlExtractor := NewHTMLParse()

	u, _ := url.Parse("https://bing.com/v1/")
	urls, err := urlExtractor.ExtractURLs(*u, []byte(siteWithValidHTML))
	if err != nil {
		t.Fatalf("extracting links from the valid site err: %s", err)
	}
	expectedURLs := []string{
		"https://bing.com/v1/test",
		"https://bing.com/test2",
		"https://google.com/test",
	}
	if len(urls) != len(expectedURLs) {
		t.Fatalf("invalid number of extracted urls from valid site, got: %d, expected: %d. URLs: %v",
			len(urls), len(expectedURLs), urls)
	}
	for i, u := range urls {
		if u.String() != expectedURLs[i] {
			t.Errorf("got: %s, expected: %s", u.String(), expectedURLs[i])
		}
	}
}
