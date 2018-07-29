package driver

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavo/filter"
)

// FetchFunc is a function that is invoked for each document found in a fetch
// operation.
//
// The fetch operation is ended if it returns false or a non-nil error.
type FetchFunc func(*document.Document) (bool, error)

// Fetch is a request to retrieve documents from the store.
type Fetch struct {
	operation

	Each   FetchFunc
	Filter *filter.Filter
}

// ExecuteInReadTx executes this operation within the context of tx.
func (o *Fetch) ExecuteInReadTx(ctx context.Context, tx ReadTx) {
	tx.Fetch(ctx, o)
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *Fetch) ExecuteInWriteTx(ctx context.Context, tx WriteTx) {
	o.ExecuteInReadTx(ctx, tx)
}
