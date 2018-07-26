package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// opLoad is an operation that loads a document by its ID.
type opLoad struct {
	DocumentID string
	Result     *Document
}

// View executes the operation.
func (op *opLoad) View(ctx context.Context, ns string, tx *bolt.Tx) (bool, error) {
	b, ok, err := getBuckets(ns, tx)
	if !ok || err != nil {
		return false, err
	}

	md, ok, err := tryGetMetaData(b, op.DocumentID)
	if !ok || err != nil {
		return false, err
	}

	c, err := getContent(b, op.DocumentID)
	if err != nil {
		return false, err
	}

	op.Result = &Document{md, c}

	return true, nil
}
