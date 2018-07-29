package protavo

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
)

// DB is a protocol buffers based document store.
type DB struct {
	ns string
	d  driver.Driver
}

// NewDB returns a new DB that uses the given driver.
func NewDB(d driver.Driver) (*DB, error) {
	return &DB{"", d}, nil
}

func (db *DB) Load(ctx context.Context, id string) (*Document, bool, error) {
	return db.LoadWhere(
		ctx,
		HasID(id),
	)
}

func (db *DB) LoadMany(ctx context.Context, ids ...string) ([]*Document, error) {
	return db.LoadManyWhere(
		ctx,
		HasID(ids...),
	)
}

func (db *DB) LoadByUniqueKey(ctx context.Context, u string) (*Document, bool, error) {
	return db.LoadWhere(
		ctx,
		HasUniqueKey(u),
	)
}

func (db *DB) LoadWhere(ctx context.Context, f ...filter.Condition) (*Document, bool, error) {
	var doc *Document

	return doc, doc != nil, db.Read(
		ctx,
		FetchWhere(
			func(d *document.Document) (bool, error) {
				doc = d
				return false, nil
			},
			f...,
		),
	)
}

func (db *DB) LoadManyWhere(ctx context.Context, f ...filter.Condition) ([]*Document, error) {
	var docs []*Document

	return docs, db.Read(
		ctx,
		FetchWhere(
			func(d *document.Document) (bool, error) {
				docs = append(docs, d)
				return true, nil
			},
			f...,
		),
	)
}

func (db *DB) FetchAll(ctx context.Context, fn driver.FetchFunc) error {
	return db.Read(
		ctx,
		FetchAll(fn),
	)
}

func (db *DB) FetchWhere(ctx context.Context, fn driver.FetchFunc, f ...filter.Condition) error {
	return db.Read(
		ctx,
		FetchWhere(fn, f...),
	)
}

func (db *DB) Save(ctx context.Context, docs ...*Document) error {
	ops := make([]driver.WriteTxOp, len(docs))

	for i, doc := range docs {
		ops[i] = Save(doc)
	}

	return db.Write(ctx, ops...)
}

func (db *DB) ForceSave(ctx context.Context, docs ...*Document) error {
	ops := make([]driver.WriteTxOp, len(docs))

	for i, doc := range docs {
		ops[i] = ForceSave(doc)
	}

	return db.Write(ctx, ops...)
}

func (db *DB) Delete(ctx context.Context, docs ...*Document) error {
	ops := make([]driver.WriteTxOp, len(docs))

	for i, doc := range docs {
		ops[i] = Delete(doc)
	}

	return db.Write(ctx, ops...)
}

func (db *DB) ForceDelete(ctx context.Context, docs ...*Document) ([]string, error) {
	ids := make([]string, len(docs))
	for i, doc := range docs {
		ids[i] = doc.ID
	}

	return db.DeleteByID(ctx, ids...)
}

func (db *DB) DeleteByID(ctx context.Context, ids ...string) ([]string, error) {
	return db.DeleteWhere(
		ctx,
		HasID(ids...),
	)
}

func (db *DB) DeleteWhere(ctx context.Context, f ...filter.Condition) ([]string, error) {
	op := DeleteWhere(f...)

	if err := db.Write(ctx, op); err != nil {
		return nil, err
	}

	return op.Result.Get()
}

// Namespace returns a DB that operates on a sub-namespace of this DB.
func (db *DB) Namespace(ns string) *DB {
	if db.ns != "" {
		ns = db.ns + "." + ns
	}

	return &DB{
		ns,
		driver.NoOpCloser{
			Driver: db.d,
		},
	}
}

// DeleteNamespace unconditionally deletes the namespace and all documents
// within it.
func (db *DB) DeleteNamespace(ctx context.Context) error {
	return db.Write(
		ctx,
		&driver.DeleteNamespace{},
	)
}

// Read atomically executes a set of read operations.
func (db *DB) Read(
	ctx context.Context,
	ops ...driver.ReadTxOp,
) error {
	tx, err := db.d.BeginRead(ctx, db.ns)
	if err != nil {
		return err
	}
	defer tx.Close()

	for _, op := range ops {
		if err := op.ExecuteInReadTx(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

// Write atomically executes a set of read/write operations.
func (db *DB) Write(
	ctx context.Context,
	ops ...driver.WriteTxOp,
) error {
	tx, err := db.d.BeginWrite(ctx, db.ns)
	if err != nil {
		return err
	}
	defer tx.Close()

	for _, op := range ops {
		if err := op.ExecuteInWriteTx(ctx, tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Close closes the DB and the underlying driver, freeing resources and
// preventing any further operations.
func (db *DB) Close() error {
	return db.d.Close()
}
