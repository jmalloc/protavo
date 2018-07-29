package protavobolt

import (
	"github.com/jmalloc/protavo/src/protavo/filter"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// isFilterSatisfiedByRecord checks if a record matches a filter.
func isFilterSatisfiedByRecord(
	c filter.Condition,
	id string,
	rec *database.Record,
) bool {
	ok, err := c.Accept(&recordMatcher{id, rec})
	if err != nil {
		panic(err)
	}

	return ok
}

// recordMatcher is a filter.Visitor that matches a filter against a document
// record.
type recordMatcher struct {
	id  string
	rec *database.Record
}

func (m *recordMatcher) IsOneOf(c *filter.IsOneOf) (bool, error) {
	_, ok := c.Values[m.id]
	return ok, nil
}

func (m *recordMatcher) HasUniqueKeyIn(c *filter.HasUniqueKeyIn) (bool, error) {
	// iterate the smaller of the two sets
	if len(c.Values) <= len(m.rec.Keys) {
		for k := range c.Values {
			t, ok := m.rec.Keys[k]
			if ok && t == database.UniqueKeyType {
				return true, nil
			}
		}
	} else {
		for k, t := range m.rec.Keys {
			if t == database.UniqueKeyType {
				if _, ok := c.Values[k]; ok {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (m *recordMatcher) HasKeys(c *filter.HasKeys) (bool, error) {
	for k := range c.Values {
		if _, ok := m.rec.Keys[k]; !ok {
			return false, nil
		}
	}

	return true, nil
}
