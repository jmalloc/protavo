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
	Each   FetchFunc
	Filter *filter.Filter
	Result *Result
}

// ExecuteInReadTx executes this operation within the context of tx.
func (o *Fetch) ExecuteInReadTx(ctx context.Context, tx ReadTx) error {
	tx.Fetch(ctx, o)
	return o.Result.Err
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *Fetch) ExecuteInWriteTx(ctx context.Context, tx WriteTx) error {
	return o.ExecuteInReadTx(ctx, tx)
}
