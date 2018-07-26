package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// driver is a low-level interface for performing database operations.
type driver interface {
	View(ctx context.Context, op viewOp) (bool, error)
	Update(ctx context.Context, op updateOp) error
	Close() error
}

// viewOp is a read-only operation.
type viewOp interface {
	View(context.Context, *bolt.Tx) (bool, error)
}

// updateOp executes is a read/write operation.
type updateOp interface {
	Update(context.Context, *bolt.Tx) error
}
