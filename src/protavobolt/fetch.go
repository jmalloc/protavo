package protavobolt

import (
	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// executeFetch calls fn for each document that matches f.
func executeFetch(
	tx *bolt.Tx,
	ns string,
	f *filter.Filter,
	fn driver.FetchFunc,
) error {
	s, ok, err := database.OpenStore(tx, ns)
	if !ok || err != nil {
		return err
	}

	return selectStrategy(s, f).Fetch(fn)
}

// applyFetch executes the side-effects of a fetch operation.
func applyFetch(
	s *database.Store,
	id string,
	rec *database.Record,
	fn driver.FetchFunc,
) (bool, error) {
	c, err := s.GetContent(id)
	if err != nil {
		return false, err
	}

	doc, err := newDocument(id, rec, c)
	if err != nil {
		return false, err
	}

	return fn(doc)
}

// Fetch is the implementation of the "no-op" strategy for fetching.
func (*noop) Fetch(fn driver.FetchFunc) error {
	return nil
}

// Fetch is the implementation of the "scan records" strategy for fetching.
func (qs *scanRecords) Fetch(fn driver.FetchFunc) error {
	cur := qs.store.Records.Cursor()

	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		rec, err := database.UnmarshalRecord(v)
		if err != nil {
			return err
		}

		id := string(k)

		if !isFilterSatisfiedByRecord(qs.filter, id, rec) {
			continue
		}

		ok, err := applyFetch(qs.store, id, rec, fn)
		if !ok || err != nil {
			return err
		}
	}

	return nil
}

// Fetch is the implementation of the "use document ID first" strategy for
// fetching.
func (qs *useIDFirst) Fetch(fn driver.FetchFunc) error {
	for id := range qs.conds.ExtractIsOneOf().Values {
		rec, exists, err := qs.store.TryGetRecord(id)
		if err != nil {
			return err
		}

		if !exists {
			continue
		}

		if !qs.conds.AreSatisfiedBy(id, rec) {
			continue
		}

		ok, err := applyFetch(qs.store, id, rec, fn)
		if !ok || err != nil {
			return err
		}
	}

	return nil
}

// Fetch is the implementation of the "use unique key first" strategy for
// fetching.
func (qs *useUniqueKeyFirst) Fetch(fn driver.FetchFunc) error {
	for key := range qs.conds.ExtractHasUniqueKeyIn().Values {
		k, err := qs.store.GetKey(key)
		if err != nil {
			return err
		}

		id, ok := k.GetUniqueDocumentID()
		if !ok {
			continue
		}

		rec, err := qs.store.GetRecord(id)
		if err != nil {
			return err
		}

		if !qs.conds.AreSatisfiedBy(id, rec) {
			continue
		}

		ok, err = applyFetch(qs.store, id, rec, fn)
		if !ok || err != nil {
			return err
		}
	}

	return nil
}

// Fetch is the implementation of the "use keys first" strategy for fetching.
func (qs *useKeysFirst) Fetch(fn driver.FetchFunc) error {
	ids, err := qs.findDocumentIDs()
	if err != nil {
		return err
	}

	for id := range ids {
		rec, err := qs.store.GetRecord(id)
		if err != nil {
			return err
		}

		if !qs.conds.AreSatisfiedBy(id, rec) {
			continue
		}

		ok, err := applyFetch(qs.store, id, rec, fn)
		if !ok || err != nil {
			return err
		}
	}

	return nil
}
