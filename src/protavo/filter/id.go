package filter

import "github.com/jmalloc/protavo/src/protavo/document"

// IsOneOf is a condition that matches documents with IDs in a specific set.
type IsOneOf struct {
	Values Set
}

// IsSatisfiedBy returns true if doc meets this condition.
func (c *IsOneOf) IsSatisfiedBy(doc *document.Document) bool {
	_, ok := c.Values[doc.ID]
	return ok
}

// Accept calls v.IsOneOf(c).
func (c *IsOneOf) Accept(v Visitor) (bool, error) {
	return v.IsOneOf(c)
}
