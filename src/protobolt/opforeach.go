package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// opForEach is an operation that iterates through all documents.
type opForEach struct {
	Fn func(*Document) error
}

func (op *opForEach) View(ctx context.Context, tx *bolt.Tx) (bool, error) {
	b, ok, err := getBuckets(tx)
	if !ok || err != nil {
		return false, err
	}

	c := b.MetaData.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		md, err := unmarshalMetaData(v)
		if err != nil {
			return false, err
		}

		c, err := getContent(b, md.Id)
		if err != nil {
			return false, err
		}

		if err := op.Fn(&Document{md, c}); err != nil {
			return false, err
		}
	}

	return true, nil
}

// opForEachMatch is an operation that iterates through all documents that
// match a set of keys.
type opForEachMatch struct {
	Fn         func(*Document) error
	Constraint []string
}

func (op *opForEachMatch) View(ctx context.Context, tx *bolt.Tx) (bool, error) {
	b, ok, err := getBuckets(tx)
	if !ok || err != nil {
		return false, err
	}

	docIDs := map[string]struct{}{}

	for _, key := range op.Constraint {
		kd, err := getKeyData(b, key)
		if err != nil {
			return false, err
		}

		for docID := range kd.Documents {
			docIDs[docID] = struct{}{}
		}
	}

	for docID := range docIDs {
		md, err := getMetaData(b, docID)
		if err != nil {
			return false, err
		}

		c, err := getContent(b, docID)
		if err != nil {
			return false, err
		}

		if err := op.Fn(&Document{md, c}); err != nil {
			return false, err
		}
	}

	return true, nil
}
