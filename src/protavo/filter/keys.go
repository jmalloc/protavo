package filter

import "github.com/jmalloc/protavo/src/protavo/document"

// MatchAllKeys is a condition that matches documents that have all of a given
// set of keys.
type MatchAllKeys struct {
	Keys []string
}

// Match returns true if doc meets this condition.
func (c *MatchAllKeys) Match(doc *document.Document) bool {
outer:
	for _, required := range c.Keys {
		for present := range doc.Keys {
			if present == required {
				continue outer
			}
		}

		return false
	}

	return true
}

// Accept calls v.MatchAllKeys(c).
func (c *MatchAllKeys) Accept(v Visitor) error {
	return v.MatchAllKeys(c)
}

// MatchUniqueKey is a condition that matches documents with unique keys in
// given set.
type MatchUniqueKey struct {
	Keys []string
}

// Match returns true if doc meets this condition.
func (c *MatchUniqueKey) Match(doc *document.Document) bool {
	for _, present := range doc.UniqueKeys() {
		for _, required := range c.Keys {
			if present == required {
				return true
			}
		}
	}

	return false
}

// Accept calls v.MatchUniqueKey(c).
func (c *MatchUniqueKey) Accept(v Visitor) error {
	return v.MatchUniqueKey(c)
}
