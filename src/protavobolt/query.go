package protavobolt

import (
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

type plan interface {
	Fetch(fn driver.FetchFunc) error
	DeleteWhere() error
}

// planQuery returns the plan used to execute an operation that applies
// to documents matching f.
func planQuery(
	s *database.Store,
	f *filter.Filter,
) (plan, error) {
	// the nil filter matches everything
	if f == nil {
		return &scanRecords{s, nil}, nil
	}

	// the empty filter matches nothing
	if len(f.Conditions) == 0 {
		return &noop{}, nil
	}

	// otherwise the best strategy we have so far is to scan all records and match
	// the filters
	return &scanRecords{s, f}, nil
}

// noop is a plan that does nothing.
// It is used when a filter is not able to match any documents.
type noop struct{}

// scanRecords is a plan that scans all records.
type scanRecords struct {
	store  *database.Store
	filter *filter.Filter
}
