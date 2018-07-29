package protavobolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/driver"
)

// ExclusiveDriver is an implementation of protavo.Driver backed by a BoltDB
// database that is held open for the life-time of the driver.
type ExclusiveDriver struct {
	DB *bolt.DB
}

// BeginRead starts a new read-only transaction.
func (d *ExclusiveDriver) BeginRead(
	ctx context.Context,
	ns string,
) (driver.ReadTx, error) {
	tx, err := d.DB.Begin(false)
	if err != nil {
		return nil, err
	}

	return &readTx{ns, tx}, nil
}

// BeginWrite starts a new read/write transaction.
func (d *ExclusiveDriver) BeginWrite(
	ctx context.Context,
	ns string,
) (driver.WriteTx, error) {
	tx, err := d.DB.Begin(true)
	if err != nil {
		return nil, err
	}

	return &writeTx{
		readTx: readTx{ns, tx},
	}, nil
}

// Close closes the driver, freeing any resources and preventing further
// operations.
func (d *ExclusiveDriver) Close() error {
	return d.DB.Close()
}
