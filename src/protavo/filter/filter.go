package filter

import "github.com/jmalloc/protavo/src/protavo/document"

// Condition is a predicate that checks if a document meets a certain criteria.
type Condition interface {
	Match(*document.Document) bool
	Accept(Visitor) error
}

// Filter is a collection of zero or more conditions.
// A filter with no conditions does not match any documents.
type Filter struct {
	Conditions []Condition
}

// New returns a new filter with the given conditions.
func New(conds []Condition) *Filter {
	return &Filter{conds}
}

// Match returns true if doc meets this condition.
func (f *Filter) Match(doc *document.Document) bool {
	if f == nil {
		// the nil filter matches everything
		return true
	} else if len(f.Conditions) == 0 {
		// the empty filter matches nothing
		return false
	}

	for _, c := range f.Conditions {
		if !c.Match(doc) {
			return false
		}
	}

	return true
}
