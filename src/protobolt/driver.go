package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// Driver is a low-level interface for performing database operations.
type Driver interface {
	View(ctx context.Context, op ViewOp) (bool, error)
	Update(ctx context.Context, op UpdateOp) error
	Close() error
}

// ViewOp is a read-only operation.
type ViewOp interface {
	View(context.Context, *bolt.Tx) (bool, error)
}

// UpdateOp executes is a read/write operation.
type UpdateOp interface {
	Update(context.Context, *bolt.Tx) error
}
