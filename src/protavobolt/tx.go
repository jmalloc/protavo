package protavobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// readTx is a BoltDB implementation of protavo.ReadTx.
type readTx struct {
	ns string
	tx *bolt.Tx
}

func (tx *readTx) Fetch(_ context.Context, op *driver.Fetch) {
	op.Result = &driver.Result{
		Err: executeFetch(
			tx.tx,
			tx.ns,
			op.Filter,
			op.Each,
		),
	}
}

func (tx *readTx) Close() error {
	return tx.tx.Rollback()
}

// readTx is a BoltDB implementation of protavo.WriteTx.
type writeTx struct {
	readTx
}

func (tx *writeTx) Save(_ context.Context, op *driver.Save) {
	op.Result = &driver.Result{
		Err: executeSave(
			tx.tx,
			tx.ns,
			op.Document,
			op.Force,
		),
	}
}

func (tx *writeTx) Delete(_ context.Context, op *driver.Delete) {
	op.Result = &driver.Result{
		Err: executeDelete(
			tx.tx,
			tx.ns,
			op.Document,
		),
	}
}

func (tx *writeTx) DeleteWhere(_ context.Context, op *driver.DeleteWhere) {
	ids, err := executeDeleteWhere(
		tx.tx,
		tx.ns,
		op.Filter,
	)

	op.Result = &driver.DeleteWhereResult{
		DocumentIDs: ids,
		Err:         err,
	}
}

func (tx *writeTx) DeleteNamespace(_ context.Context, op *driver.DeleteNamespace) {
	op.Result = &driver.Result{
		Err: database.DeleteStore(tx.tx, tx.ns),
	}
}

func (tx *writeTx) Commit() error {
	return tx.tx.Commit()
}
