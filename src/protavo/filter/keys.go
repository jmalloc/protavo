package filter

import "github.com/jmalloc/protavo/src/protavo/document"

// HasKeys is a condition that matches documents that have all of a given
// set of keys.
type HasKeys struct {
	Values Set
}

// IsSatisfiedBy returns true if doc meets this condition.
func (c *HasKeys) IsSatisfiedBy(doc *document.Document) bool {
	for k := range c.Values {
		if _, ok := doc.Keys[k]; !ok {
			return false
		}
	}

	return true
}

// Accept calls v.HasKeys(c).
func (c *HasKeys) Accept(v Visitor) (bool, error) {
	return v.HasKeys(c)
}

// HasUniqueKeyIn is a condition that matches documents with unique keys in
// given set.
type HasUniqueKeyIn struct {
	Values Set
}

// IsSatisfiedBy returns true if doc meets this condition.
func (c *HasUniqueKeyIn) IsSatisfiedBy(doc *document.Document) bool {
	// iterate the smaller of the two sets
	if len(c.Values) <= len(doc.Keys) {
		for k := range c.Values {
			t, ok := doc.Keys[k]
			if ok && t == document.UniqueKey {
				return true
			}
		}
	} else {
		for k, t := range doc.Keys {
			if t == document.UniqueKey {
				if _, ok := c.Values[k]; ok {
					return true
				}
			}
		}
	}

	return false
}

// Accept calls v.HasUniqueKeyIn(c).
func (c *HasUniqueKeyIn) Accept(v Visitor) (bool, error) {
	return v.HasUniqueKeyIn(c)
}
