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

	p, err := planQuery(s, f)
	if err != nil {
		return err
	}

	return p.Fetch(fn)
}

func (*noop) Fetch(fn driver.FetchFunc) error {
	return nil
}

func (p *scanRecords) Fetch(fn driver.FetchFunc) error {
	cur := p.store.Records.Cursor()

	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		rec, err := database.UnmarshalRecord(v)
		if err != nil {
			return err
		}

		id := string(k)

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

		ok, err := fn(doc)
		if !ok || err != nil {
			return err
		}
	}

	return nil
}
