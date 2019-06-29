package sitemap

import (
	"bytes"
)

func generatePlaintext(entries []Entry) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	for _, entry := range entries {
		buf.WriteString(entry.Location.String() + "\n")
	}
	return buf.Bytes(), nil
}
