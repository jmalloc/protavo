package protavobolt

import (
	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// executeDeleteWhere deletes documents that match f, regardless of whether
// their revisions match the currently persisted revisions.
func executeDeleteWhere(
	tx *bolt.Tx,
	ns string,
	f *filter.Filter,
	fn driver.DeleteWhereFunc,
) error {
	s, ok, err := database.OpenStore(tx, ns)
	if !ok || err != nil {
		return err
	}

	return selectStrategy(s, f).DeleteWhere(fn)
}

// applyDelete executes the side-effects of a delete-where operation.
func applyDelete(
	s *database.Store,
	id string,
	rec *database.Record,
	deleteRec bool,
	fn driver.DeleteWhereFunc,
) error {
	if deleteRec {
		if err := s.DeleteRecord(id); err != nil {
			return err
		}
	}

	if err := s.DeleteContent(id); err != nil {
		return err
	}

	if err := s.UpdateKeys(id, rec.Keys, nil); err != nil {
		return err
	}

	if fn != nil {
		return fn(id)
	}

	return nil
}

// DeleteWhere is the implementation of the "no-op" strategy for deleting.
func (*noop) DeleteWhere(driver.DeleteWhereFunc) error {
	return nil
}

// DeleteWhere is the implementation of the "scan records" strategy for deleting.
func (qs *scanRecords) DeleteWhere(fn driver.DeleteWhereFunc) error {
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

		// use the cursor to delete the record so that we don't invalidate it
		if err := cur.Delete(); err != nil {
			return err
		}

		if err := applyDelete(qs.store, id, rec, false, fn); err != nil {
			return err
		}
	}

	return nil
}

// DeleteWhere is the implementation of the "use document ID first" strategy for
// deleting.
func (qs *useIDFirst) DeleteWhere(fn driver.DeleteWhereFunc) error {
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

		if err := applyDelete(qs.store, id, rec, true, fn); err != nil {
			return err
		}
	}

	return nil
}

// DeleteWhere is the implementation of the "use unique key first" strategy for
// deleting.
func (qs *useUniqueKeyFirst) DeleteWhere(fn driver.DeleteWhereFunc) error {
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

		if err := applyDelete(qs.store, id, rec, true, fn); err != nil {
			return err
		}
	}

	return nil
}

// DeleteWhere is the implementation of the "use keys first" strategy for
// deleting.
func (qs *useKeysFirst) DeleteWhere(fn driver.DeleteWhereFunc) error {
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

		if err := applyDelete(qs.store, id, rec, true, fn); err != nil {
			return err
		}
	}

	return nil
}
