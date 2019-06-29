package sitemap

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type xmlURLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`

	URLs []xmlURL `xml:"url"`
}

type xmlURL struct {
	Location        string    `xml:"loc"`
	ChangeFrequency Frequency `xml:"changefreq,omitempty"`
	Priority        float64   `xml:"priority,omitempty"`
}

func generateXML(entries []Entry) ([]byte, error) {
	urls := make([]xmlURL, 0, len(entries))
	for _, entry := range entries {
		urls = append(urls, xmlURL{
			Location:        entry.Location.String(),
			ChangeFrequency: entry.ChangeFrequency,
			Priority:        entry.Priority,
		})
	}
	root := xmlURLSet{
		XMLNS: xmlns,
		URLs:  urls,
	}
	data, err := xml.Marshal(root)
	if err != nil {
		return nil, errors.Wrapf(err, "marshaling")
	}
	return data, nil
}
