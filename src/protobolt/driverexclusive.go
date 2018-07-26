package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// ExclusiveDriver is an implementation of Driver that uses a single *bolt.DB over
// its entire lifetime, preventing other processes from opening the database.
type ExclusiveDriver struct {
	Database *bolt.DB
}

// View executes a read-only operation.
func (d *ExclusiveDriver) View(ctx context.Context, op ViewOp) (bool, error) {
	var ok bool

	return ok, d.Database.View(func(tx *bolt.Tx) error {
		var err error
		ok, err = op.View(ctx, tx)
		return err
	})
}

// Update executes a read/write operation.
func (d *ExclusiveDriver) Update(ctx context.Context, op UpdateOp) error {
	return d.Database.Update(func(tx *bolt.Tx) error {
		return op.Update(ctx, tx)
	})
}

// Close closes the database.
func (d *ExclusiveDriver) Close() error {
	return d.Database.Close()
}
