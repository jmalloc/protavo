package filter

import "github.com/jmalloc/protavo/src/protavo/document"

// MatchNothing is a condition that never matches any documents.
type MatchNothing struct {
}

// Match returns true if doc meets this condition.
func (c *MatchNothing) Match(*document.Document) bool {
	return false
}

// Accept calls v.MatchNothing(c).
func (c *MatchNothing) Accept(v Visitor) error {
	return v.MatchNothing(c)
}
