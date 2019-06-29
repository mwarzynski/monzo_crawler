package url_extractor

import (
	"net/url"
	"strings"
)

func normalizeURL(u url.URL) url.URL {
	u.Fragment = ""
	u.Path = strings.Replace(u.Path, "//", "/", -1)
	u.RawPath = strings.Replace(u.RawPath, "//", "/", -1)
	u.ForceQuery = false
	return u
}
