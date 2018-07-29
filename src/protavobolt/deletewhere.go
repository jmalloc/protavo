package protavobolt

import (
	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/filter"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// executeDeleteWhere deletes documents that match f, regardless of whether
// their revisions match the currently persisted revisions.
func executeDeleteWhere(
	tx *bolt.Tx,
	ns string,
	f *filter.Filter,
) ([]string, error) {
	panic("ni")
}

func (*noop) DeleteWhere() error {
	return nil
}

func (p *scanRecords) DeleteWhere() error {
	cur := p.store.Records.Cursor()

	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		rec, err := database.UnmarshalRecord(v)
		if err != nil {
			return err
		}

		id := string(k)

		// we have to construct the whole document just to match the filter
		if p.filter != nil {
			c, err := p.store.GetContent(id)
			if err != nil {
				return err
			}

			doc, err := newDocument(id, rec, c)
			if err != nil {
				return err
			}

			if !p.filter.Match(doc) {
				continue
			}
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
	}

	return nil
}
