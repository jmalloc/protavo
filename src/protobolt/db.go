package protobolt

import (
	"context"
	"os"

	bolt "github.com/coreos/bbolt"
)

// DB is a simple protocol buffers based document store.
type DB struct {
	Driver Driver
}

// Open creates an opens a new database at the given path.
//
// If shared is true, the database is only opened when performing an operation,
// allow the database to be shared between multiple processes.
func Open(
	path string,
	shared bool,
	mode os.FileMode,
	opts *bolt.Options,
) (*DB, error) {
	if shared {
		return &DB{
			&SharedDriver{
				Path:    path,
				Mode:    mode,
				Options: opts,
			},
		}, nil
	}

	db, err := bolt.Open(path, mode, opts)
	if err != nil {
		return nil, err
	}

	return &DB{
		&ExclusiveDriver{
			Database: db,
		},
	}, nil
}

// Load returns the document with the given ID.
//
// It returns false if the document does not exist.
func (db *DB) Load(ctx context.Context, id string) (*Document, bool, error) {
	op := &opLoad{
		DocumentID: id,
	}

	ok, err := db.Driver.View(ctx, op)

	return op.Document, ok, err
}

// Save persists a document to the store, creating it if it does not already
// exist.
//
// doc.MetaData().Version must match the currently persisted version of the
// document, which is zero for non-existent documents. Otherwise; an
// OptimisticLockError is returned.
//
// It returns the saved version of the document, with the new version number.
func (db *DB) Save(ctx context.Context, doc *Document) (*Document, error) {
	saved, err := db.SaveMany(ctx, doc)
	if err != nil {
		return nil, err
	}

	return saved[0], nil
}

// SaveMany atomically persists one or more documents to the store, creating
// any documents that do not already exists.
//
// For each document, doc.MetaData().Version must match the currently persisted
// version of that document, which is zero for non-existent documents.
// Otherwise; an OptimisticLockError is returned.
//
// It returns the saved versions of the documents, in the same order as they
// are provided.
func (db *DB) SaveMany(ctx context.Context, docs ...*Document) ([]*Document, error) {
	if len(docs) == 0 {
		return nil, nil
	}

	op := &opSave{
		Documents: docs,
	}

	return op.SavedDocuments, db.Driver.Update(ctx, op)
}

// Delete atomically removes one or more documents from the store.
//
// For each document, doc.MetaData().Version must match the currently persisted
// version of that document, which is zero for non-existent documents.
// Otherwise; an OptimisticLockError is returned. It is not an error to delete
// a non-existent document provided that the given version is zero.
func (db *DB) Delete(ctx context.Context, docs ...*Document) error {
	if len(docs) == 0 {
		return nil
	}

	return db.Driver.Update(ctx, &opDelete{docs})
}

// Find returns the the document that has the given unique key.
//
// u must refer to a unique key. If no such key exists, or the key is not
// unique, false is returned.
//
// filter is additional set of keys (either unique or shared) that the document
// must have in order to match.
func (db *DB) Find(ctx context.Context, uniq string, keys ...string) (*Document, bool, error) {
	op := &opFind{
		UniqueKey:  uniq,
		Constraint: keys,
	}

	ok, err := db.Driver.View(ctx, op)

	return op.Document, ok, err
}

// FindMany returns the documents that have all of the keys in the given
// filter.
func (db *DB) FindMany(ctx context.Context, keys ...string) ([]*Document, error) {
	var docs []*Document

	fn := func(doc *Document) error {
		docs = append(docs, doc)
		return nil
	}

	return docs, db.ForEach(ctx, fn, keys...)
}

// ForEach calls fn for each document that has all of the keuys in the given
// filter.
//
// If fn returns an error, iteration stops and that error is returned.
func (db *DB) ForEach(
	ctx context.Context,
	fn func(*Document) error,
	keys ...string,
) error {
	var op ViewOp

	if len(keys) == 0 {
		op = &opForEach{fn}
	} else {
		op = &opForEachMatch{fn, keys}
	}

	_, err := db.Driver.View(ctx, op)
	return err
}

// Close closes the database.
func (db *DB) Close() error {
	return db.Driver.Close()
}
