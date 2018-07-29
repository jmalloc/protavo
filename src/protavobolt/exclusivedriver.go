package protavobolt

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/protavo/src/protavo/driver"
)

// ExclusiveDriver is an implementation of protavo.Driver backed by a BoltDB
// database that is held open for the life-time of the driver.
type ExclusiveDriver struct {
	DB      *bolt.DB
	onClose func() error
}

// OpenExclusive returns a new ExclusiveDriver.
func OpenExclusive(file string, mode os.FileMode, opts *bolt.Options) (*ExclusiveDriver, error) {
	db, err := bolt.Open(file, mode, opts)
	if err != nil {
		return nil, err
	}

	return &ExclusiveDriver{db, nil}, nil
}

// OpenTemp opens a new driver that operates on a temporary file.
func OpenTemp(mode os.FileMode, opts *bolt.Options) (*ExclusiveDriver, error) {
	dir, err := ioutil.TempDir("", "protavobolt-")
	if err != nil {
		return nil, err
	}

	file := path.Join(dir, "bolt.db")
	db, err := bolt.Open(file, mode, opts)
	if err != nil {
		return nil, err
	}

	return &ExclusiveDriver{
		db,
		func() error {
			return os.RemoveAll(dir)
		},
	}, nil
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
	err := d.DB.Close()

	if d.onClose != nil {
		e := d.onClose()

		// we always call onClose when it is present, even if Close() returned an
		// error, but always favour returning the close error as it probably has more
		// important information.
		if err == nil {
			err = e
		}
	}

	return err
}
