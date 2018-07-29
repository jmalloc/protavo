package protavobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// readTx is a BoltDB implementation of protavo.ReadTx.
//
// TODO(jmalloc): put a 'now' timestamp in the tx so operations in the same tx
// share the same timestamp.
type readTx struct {
	ns string
	tx *bolt.Tx
}

func (tx *readTx) Fetch(_ context.Context, op *driver.Fetch) {
	op.MarkExecuted(
		executeFetch(
			tx.tx,
			tx.ns,
			op.Filter,
			op.Each,
		),
	)
}

func (tx *readTx) Close() error {
	return tx.tx.Rollback()
}

// readTx is a BoltDB implementation of protavo.WriteTx.
type writeTx struct {
	readTx
}

func (tx *writeTx) Save(_ context.Context, op *driver.Save) {
	op.MarkExecuted(
		executeSave(
			tx.tx,
			tx.ns,
			op.Document,
			op.Force,
		),
	)
}

func (tx *writeTx) Delete(_ context.Context, op *driver.Delete) {
	op.MarkExecuted(
		executeDelete(
			tx.tx,
			tx.ns,
			op.Document,
		),
	)
}

func (tx *writeTx) DeleteWhere(_ context.Context, op *driver.DeleteWhere) {
	op.MarkExecuted(
		executeDeleteWhere(
			tx.tx,
			tx.ns,
			op.Filter,
			op.Each,
		),
	)
}

func (tx *writeTx) DeleteNamespace(_ context.Context, op *driver.DeleteNamespace) {
	op.MarkExecuted(
		database.DeleteStore(
			tx.tx,
			tx.ns,
		),
	)
}

func (tx *writeTx) Commit() error {
	return tx.tx.Commit()
}
