package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protobolt/src/protobolt/internal/types"
)

// opFind is an operation that finds a single document by a unique key.
type opFind struct {
	UniqueKey  string
	Constraint []string
	Document   *Document
}

func (op *opFind) View(ctx context.Context, tx *bolt.Tx) (bool, error) {
	b, ok, err := getBuckets(tx)
	if !ok || err != nil {
		return false, err
	}

	kd, err := getKeyData(b, op.UniqueKey)
	if err != nil {
		return false, err
	}

	// if the key is not unique, it does not identify a single document
	if kd.Type != types.KeyType_UNIQUE {
		return false, nil
	}

	for docID := range kd.Documents {
		md, err := getMetaData(b, docID)
		if err != nil {
			return false, err
		}

		if !matchesConstraint(md, op.Constraint) {
			return false, nil
		}

		c, err := getContent(b, docID)
		if err != nil {
			return false, err
		}

		op.Document = &Document{md, c}

		return true, nil
	}

	return false, nil
}
