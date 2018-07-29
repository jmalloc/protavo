package driver

import (
	"context"
)

// ReadTx is a transaction that can not modify the database.
type ReadTx interface {
	Fetch(ctx context.Context, op *Fetch)

	Close() error
}

// WriteTx is a transaction that can read and modify the database.
type WriteTx interface {
	ReadTx

	Save(ctx context.Context, op *Save)
	Delete(ctx context.Context, op *Delete)
	DeleteWhere(ctx context.Context, op *DeleteWhere)
	DeleteNamespace(ctx context.Context, op *DeleteNamespace)

	Commit() error
}

// Operation is a database operation that can be performed within a transaction.
type Operation interface {
	// Err returns the error that occurred when the operation was executed, if any.
	// It panics if the operation has not yet been executed.
	Err() error

	// MarkExecuted is called by driver implementors to indicate the result of the
	// operation.
	MarkExecuted(err error)
}

// WriteTxOp is an operation that can be performed inside transactions that
// support writing.
type WriteTxOp interface {
	Operation

	// ExecuteInWriteTx executes the operation within the context of a write
	// transaction.
	ExecuteInWriteTx(ctx context.Context, tx WriteTx)
}

// ReadTxOp is an operation that can be performed inside transactions that
// support reading.
type ReadTxOp interface {
	WriteTxOp // all read operations can be performed insite write transactions.

	// ExecuteInReadTx executes the operation within the context of a read
	// transaction.
	ExecuteInReadTx(ctx context.Context, tx ReadTx)
}

// operation provides common methods to operation implementations.
type operation struct {
	err *error
}

func (o *operation) Err() error {
	if o.err == nil {
		panic("the operation has not yet been executed")
	}

	return *o.err
}

func (o *operation) MarkExecuted(err error) {
	if o.err != nil {
		panic("the operation has already been executed")
	}

	o.err = &err
}
