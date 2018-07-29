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

	p, err := planQuery(s, f)
	if err != nil {
		return err
	}

	return p.DeleteWhere(fn)
}

func (*noop) DeleteWhere(driver.DeleteWhereFunc) error {
	return nil
}

func (p *scanRecords) DeleteWhere(fn driver.DeleteWhereFunc) error {
	cur := p.store.Records.Cursor()

	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		rec, err := database.UnmarshalRecord(v)
		if err != nil {
			return err
		}

		id := string(k)

		if !isFilterSatisfiedByRecord(p.filter, id, rec) {
			continue
		}

		// use the cursor to delete the record so that we don't invalidate it
		if err := cur.Delete(); err != nil {
			return err
		}

		if err := p.store.DeleteContent(id); err != nil {
			return err
		}

		if err := p.store.UpdateKeys(id, rec.Keys, nil); err != nil {
			return err
		}

		if fn != nil {
			if err := fn(id); err != nil {
				return err
			}
		}
	}

	return nil
}

// func (p *getByID) DeleteWhere(driver.DeleteWhereFunc) error {
// 	panic("ni")
// }

// func (p *getByUniqueKey) DeleteWhere(driver.DeleteWhereFunc) error {
// 	panic("ni")
// }

// func (p *getByKeys) DeleteWhere(driver.DeleteWhereFunc) error {
// 	panic("ni")
// }
