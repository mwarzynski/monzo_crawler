package sitemap

import "github.com/pkg/errors"

type Type string

const (
	TypeXML       Type = "xml"
	TypePlaintext Type = "plaintext"
)

type Generator struct {
	Entries []Entry
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) AddEntry(entry Entry) {
	g.Entries = append(g.Entries, entry)
}

func (g *Generator) Generate(t Type) ([]byte, error) {
	switch t {
	case TypeXML:
		return generateXML(g.Entries)
	case TypePlaintext:
		return generatePlaintext(g.Entries)
	default:
		return nil, errors.Errorf("unsupported generator type '%s'", t)
	}
}
