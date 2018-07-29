package driver

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavo/filter"
)

// Delete is a request to delete a document.
type Delete struct {
	Document *document.Document
	Result   *Result
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *Delete) ExecuteInWriteTx(ctx context.Context, tx WriteTx) error {
	tx.Delete(ctx, o)
	return o.Result.Err
}

// DeleteWhere is a request to delete one or more documents, without checking the
// current document revisions.
type DeleteWhere struct {
	Filter *filter.Filter
	Result *DeleteWhereResult
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *DeleteWhere) ExecuteInWriteTx(ctx context.Context, tx WriteTx) error {
	tx.DeleteWhere(ctx, o)
	return o.Result.Err
}

// DeleteWhereResult is the result of a DeleteWhere operation.
type DeleteWhereResult struct {
	DocumentIDs []string
	Err         error
}

// Get returns the result value and error.
// It panics if the operation has not yet been executed.
func (r *DeleteWhereResult) Get() ([]string, error) {
	if r == nil {
		panic("operation has not been executed")
	}

	return r.DocumentIDs, r.Err
}

// DeleteNamespace is a request to delete an entire namespace.
type DeleteNamespace struct {
	Documents []*document.Document
	Result    *Result
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *DeleteNamespace) ExecuteInWriteTx(ctx context.Context, tx WriteTx) error {
	tx.DeleteNamespace(ctx, o)
	return o.Result.Err
}
