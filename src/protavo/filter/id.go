package filter

import "github.com/jmalloc/protavo/src/protavo/document"

// MatchDocumentID is a condition that matches documents with IDs in a specific set.
type MatchDocumentID struct {
	DocumentIDs []string
}

// Match returns true if doc meets this condition.
func (c *MatchDocumentID) Match(doc *document.Document) bool {
	for _, id := range c.DocumentIDs {
		if doc.ID == id {
			return true
		}
	}

	return false
}

// Accept calls v.MatchDocumentID(c).
func (c *MatchDocumentID) Accept(v Visitor) error {
	return v.MatchDocumentID(c)
}
