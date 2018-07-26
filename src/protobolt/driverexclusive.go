package protobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
)

// exclusiveDriver is an implementation of Driver that uses a single *bolt.DB over
// its entire lifetime, preventing other processes from opening the database.
type exclusiveDriver struct {
	Database *bolt.DB
}

// View executes a read-only operation.
func (d *exclusiveDriver) View(ctx context.Context, op viewOp) (bool, error) {
	var ok bool

	return ok, d.Database.View(func(tx *bolt.Tx) error {
		var err error
		ok, err = op.View(ctx, tx)
		return err
	})
}

// Update executes a read/write operation.
func (d *exclusiveDriver) Update(ctx context.Context, op updateOp) error {
	return d.Database.Update(func(tx *bolt.Tx) error {
		return op.Update(ctx, tx)
	})
}

// Close closes the database.
func (d *exclusiveDriver) Close() error {
	return d.Database.Close()
}
