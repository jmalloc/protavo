package protavobolt

import (
	"errors"
	"math"

	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// A strategy encapsulates the "strategy" used to implement operations that
// locate documents using a filter.
type strategy interface {
	Fetch(fn driver.FetchFunc) error
	DeleteWhere(fn driver.DeleteWhereFunc) error
}

// noop is a query strategy that does nothing.
// It is used when a filter is not able to match any documents.
type noop struct{}

// scanRecords is a query strategy that iterates over all records and passes
// them through a filter in memory.
type scanRecords struct {
	store  *database.Store
	filter *filter.Filter
}

// useIDFirst is a query strategy that retreives records by their ID, then
// applies the remaining set of filters in-memory.
type useIDFirst struct {
	store *database.Store
	conds *conditions
}

// useUniqueKeyFirst is a query strategy that retreives records by a unique-key,
// then applies the remaining set of filters in-memory.
type useUniqueKeyFirst struct {
	store *database.Store
	conds *conditions
}

// useKeysFirst is a query strategy that finds the intersection of documents
// that have all of a specific set of keys, then applies the remaining set of
// filters in-memory.
type useKeysFirst struct {
	store *database.Store
	conds *conditions
}

// findDocumentIDs returns the IDs of the documents that have all of the
// required keys
func (qs *useKeysFirst) findDocumentIDs() (map[string]bool, error) {
	var ids map[string]bool

	for key := range qs.conds.ExtractHasKeys().Values {
		k, err := qs.store.GetKey(key)
		if err != nil {
			return nil, err
		}

		if ids == nil {
			ids = k.Documents
		} else {
			for id := range ids {
				if _, ok := k.Documents[id]; !ok {
					delete(ids, id)
				}
			}
		}

		if len(ids) == 0 {
			return nil, nil
		}
	}

	return ids, nil
}

// selectStrategy returns the plan used to execute an operation that applies to documents
// matching f.
func selectStrategy(s *database.Store, f *filter.Filter) strategy {
	f = filter.Optimize(f)

	if f == nil {
		// if there's no filter, scan everything
		return &scanRecords{s, nil}
	} else if len(f.Conditions) == 0 {
		// or if the filter matches nothing, perform a noop
		return &noop{}
	}

	conds := &conditions{}
	if _, err := f.Accept(conds); err != nil {
		panic(err)
	}

	// find which condition type is the MOST constrained
	//
	// TODO(jmalloc): the initial cost should probably be derived from the number
	// of documents in the store, likewise the number of keys for key-related
	// strategies, however I don't think this should get any more complex until
	// there are benchmarks in place.
	cheapest := math.MaxUint32
	var qs strategy = &scanRecords{s, f}

	if conds.IsOneOfCondition != nil {
		cheapest = len(conds.IsOneOfCondition.Values)
		qs = &useIDFirst{s, conds}
	}

	if conds.HasUniqueKeyInCondition != nil {
		cost := len(conds.HasUniqueKeyInCondition.Values)
		if cost < cheapest {
			cheapest = cost
			qs = &useUniqueKeyFirst{s, conds}
		}
	}

	if conds.HasKeysCondition != nil {
		cost := len(conds.HasKeysCondition.Values)
		if cost < cheapest {
			cheapest = cost
			qs = &useKeysFirst{s, conds}
		}
	}

	return qs
}

// conditions is a filter.Visitor that extracts the individual constraint types
// from a filter.
//
// It relies on the fact that filter.Optimize() currently ensures there will be
// at most one of each condition type, but this is not guaranteed going forward.
type conditions struct {
	IsOneOfCondition        *filter.IsOneOf
	HasUniqueKeyInCondition *filter.HasUniqueKeyIn
	HasKeysCondition        *filter.HasKeys
}

func (x *conditions) IsOneOf(c *filter.IsOneOf) (bool, error) {
	if x.IsOneOfCondition != nil {
		return false, errors.New(
			"conditions are expected to be flattened by filter.Optimize()",
		)
	}

	x.IsOneOfCondition = c
	return true, nil
}

func (x *conditions) HasUniqueKeyIn(c *filter.HasUniqueKeyIn) (bool, error) {
	if x.HasUniqueKeyInCondition != nil {
		return false, errors.New(
			"conditions are expected to be flattened by filter.Optimize()",
		)
	}

	x.HasUniqueKeyInCondition = c
	return true, nil
}

func (x *conditions) HasKeys(c *filter.HasKeys) (bool, error) {
	if x.HasKeysCondition != nil {
		return false, errors.New(
			"conditions are expected to be flattened by filter.Optimize()",
		)
	}

	x.HasKeysCondition = c
	return true, nil
}

// ExtractIsOneOf extracts the 'IsOneOf', clearing it from x such that future
// calls to x.AreSatisfiedBy() do not check this condition.
func (x *conditions) ExtractIsOneOf() *filter.IsOneOf {
	if x.IsOneOfCondition == nil {
		panic("x.IsOneOfCondition is nil")
	}

	c := x.IsOneOfCondition
	x.IsOneOfCondition = nil

	return c
}

// ExtractHasUniqueKeyIn extracts the 'HasUniqueKeyIn', clearing it from x such
// that future calls to x.AreSatisfiedBy() do not check this condition.
func (x *conditions) ExtractHasUniqueKeyIn() *filter.HasUniqueKeyIn {
	if x.HasUniqueKeyInCondition == nil {
		panic("x.HasUniqueKeyInCondition is nil")
	}

	c := x.HasUniqueKeyInCondition
	x.HasUniqueKeyInCondition = nil

	return c
}

// ExtractHasKeys extracts the 'HasKeys', clearing it from x such that future
// calls to x.AreSatisfiedBy() do not check this condition.
func (x *conditions) ExtractHasKeys() *filter.HasKeys {
	if x.HasKeysCondition == nil {
		panic("x.HasKeysCondition is nil")
	}

	c := x.HasKeysCondition
	x.HasKeysCondition = nil

	return c
}

// AreSatisfiedBy verifies that any of the remaining non-nil conditions on
// x are met by the given record.
func (x *conditions) AreSatisfiedBy(
	id string,
	rec *database.Record,
) bool {
	if x.IsOneOfCondition != nil &&
		!isFilterSatisfiedByRecord(x.IsOneOfCondition, id, rec) {
		return false
	}

	if x.HasUniqueKeyInCondition != nil &&
		!isFilterSatisfiedByRecord(x.HasUniqueKeyInCondition, id, rec) {
		return false
	}

	if x.HasKeysCondition != nil &&
		!isFilterSatisfiedByRecord(x.HasKeysCondition, id, rec) {
		return false
	}

	return true
}
