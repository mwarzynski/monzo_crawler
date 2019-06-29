package sitemap

import (
	"net/url"
	"testing"
)

var sitemapXML = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://google.com/test1</loc></url><url><loc>https://google.com/test2</loc></url><url><loc>https://google.com/test3</loc></url></urlset>`

func urlFromString(v string) url.URL {
	u, _ := url.Parse(v)
	return *u
}

func TestGeneratorXML(t *testing.T) {
	generator := NewGenerator()

	generator.AddEntry(Entry{Location: urlFromString("https://google.com/test1")})
	generator.AddEntry(Entry{Location: urlFromString("https://google.com/test2")})
	generator.AddEntry(Entry{Location: urlFromString("https://google.com/test3")})

	sitemap, err := generator.Generate(TypeXML)
	if err != nil {
		t.Fatalf("generating XML sitemap err: %s", err)
	}
	if string(sitemap) != sitemapXML {
		t.Errorf("invalid sitemap:\n%s", string(sitemap))
	}
}
