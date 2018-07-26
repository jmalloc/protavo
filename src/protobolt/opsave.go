package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/golang/protobuf/ptypes"
)

// opSave is an operation that atomically creates/updates one or more documents.
type opSave struct {
	Documents     []*Document
	CheckVersions bool
	Result        []*Document
}

// Update executes the operation.
func (op *opSave) Update(ctx context.Context, ns string, tx *bolt.Tx) error {
	b, err := createBuckets(ns, tx)
	if err != nil {
		return err
	}

	op.Result = make([]*Document, len(op.Documents))

	for i, doc := range op.Documents {
		doc.validate()

		after := doc.cloneMetaData()
		before, alreadyExists, err := tryGetMetaData(b, after.Id)
		if err != nil {
			return err
		}

		if op.CheckVersions {
			if err := checkVersion("save", before, after); err != nil {
				return err
			}
		}

		after.UpdatedAt = ptypes.TimestampNow()

		if alreadyExists {
			after.Version = before.Version + 1
			after.CreatedAt = before.CreatedAt
		} else {
			after.Version = 1
			after.CreatedAt = after.UpdatedAt
		}

		if err := updateKeys(b, after.Id, before, after); err != nil {
			return err
		}

		if err := putMetaData(b, after); err != nil {
			return err
		}

		if err := putContent(b, after.Id, doc.c); err != nil {
			return err
		}

		op.Result[i] = &Document{after, doc.c}
	}

	return nil
}
