package driver

import "context"

// Driver is an interface for performing operations on a data store.
type Driver interface {
	// BeginRead starts a new read-only transaction.
	BeginRead(ctx context.Context, ns string) (ReadTx, error)

	// BeginWrite starts a new read/write transaction.
	BeginWrite(ctx context.Context, ns string) (WriteTx, error)

	// Close closes the driver, freeing any resources and preventing further
	// operations.
	Close() error
}

// NoOpCloser is a Driver decoarator that prevents calls to Close() from
// reaching the underlying driver.
type NoOpCloser struct {
	Driver
}

// Close is a no-op that always returns nil.
func (NoOpCloser) Close() error {
	return nil
}
