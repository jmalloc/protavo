package filter

import "github.com/jmalloc/protavo/src/protavo/document"

// Condition is a predicate that checks if a document meets a certain criteria.
type Condition interface {
	Match(*document.Document) bool
	Accept(Visitor) (bool, error)
}

// Filter is a collection of zero or more conditions.
// A filter with no conditions does not match any documents.
type Filter struct {
	Conditions []Condition
}

// New returns a new filter with the given conditions.
func New(conds []Condition) *Filter {
	return optimize(&Filter{conds})
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

// Accept calls c.Accept(v) for each condition in f.
//
// If f is nil, Accept() returns true. If f is non-nil, but contains no matchers
// it returns false. The behavior is useful for building custom matches.
func (f *Filter) Accept(v Visitor) (bool, error) {
	if f == nil {
		return true, nil
	} else if len(f.Conditions) == 0 {
		return false, nil
	}

	for _, c := range f.Conditions {
		ok, err := c.Accept(v)
		if !ok || err != nil {
			return false, err
		}
	}

	return true, nil
}
