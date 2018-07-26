package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// opDelete is an operation that atomically deletes one or more documents.
type opDelete struct {
	Documents     []*Document
	CheckVersions bool
}

// Update executes the operation.
func (op *opDelete) Update(ctx context.Context, ns string, tx *bolt.Tx) error {
	b, ok, err := getBuckets(ns, tx)
	if !ok || err != nil {
		return err
	}

	for _, doc := range op.Documents {
		doc.validate()

		md, _, err := tryGetMetaData(b, doc.md.Id)
		if err != nil {
			return err
		}

		if op.CheckVersions {
			if err := checkVersion("delete", md, doc.md); err != nil {
				return err
			}
		}

		if err := updateKeys(b, doc.md.Id, md, nil); err != nil {
			return err
		}

		docID := []byte(doc.md.Id)

		if err := b.MetaData.Delete(docID); err != nil {
			return err
		}

		if err := b.Content.Delete(docID); err != nil {
			return err
		}
	}

	return nil
}

// opDeleteAll is an operation that atomically removes all documents from the store.
type opDeleteAll struct {
}

// Update executes the operation.
func (op *opDeleteAll) Update(ctx context.Context, ns string, tx *bolt.Tx) error {
	return deleteBuckets(ns, tx)
}
