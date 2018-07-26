package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// driver is a low-level interface for performing database operations.
type driver interface {
	View(ctx context.Context, ns string, op viewOp) (bool, error)
	Update(ctx context.Context, ns string, op updateOp) error
	Close() error
}

// viewOp is a read-only operation.
type viewOp interface {
	View(ctx context.Context, ns string, tx *bolt.Tx) (bool, error)
}

// updateOp executes is a read/write operation.
type updateOp interface {
	Update(ctx context.Context, ns string, tx *bolt.Tx) error
}
