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

// ReadTxOp is an operation that can be performed inside transactions that
// support reading.
type ReadTxOp interface {
	// All read operations can be performed inside write transactions.
	WriteTxOp

	// ExecuteInReadTx executes the operation within the context of a read
	// transaction.
	ExecuteInReadTx(ctx context.Context, tx ReadTx) error
}

// WriteTxOp is an operation that can be performed inside transactions that
// support writing.
type WriteTxOp interface {
	// ExecuteInWriteTx executes the operation within the context of a write
	// transaction.
	ExecuteInWriteTx(ctx context.Context, tx WriteTx) error
}
