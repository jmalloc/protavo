package driver

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavo/filter"
)

// Delete is a request to delete a document.
type Delete struct {
	operation

	Document *document.Document
	Result   *Result
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *Delete) ExecuteInWriteTx(ctx context.Context, tx WriteTx) {
	tx.Delete(ctx, o)
}

// DeleteWhereFunc is a function that is invoked for each document deleted in a
// delete-where operation.
//
// The delete operation is aborted if it returns a non-nil error.
type DeleteWhereFunc func(id string) error

// DeleteWhere is a request to delete one or more documents, without checking the
// current document revisions.
type DeleteWhere struct {
	operation

	Each   DeleteWhereFunc
	Filter *filter.Filter
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *DeleteWhere) ExecuteInWriteTx(ctx context.Context, tx WriteTx) {
	tx.DeleteWhere(ctx, o)
}

// DeleteNamespace is a request to delete an entire namespace.
type DeleteNamespace struct {
	operation

	Documents []*document.Document
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *DeleteNamespace) ExecuteInWriteTx(ctx context.Context, tx WriteTx) {
	tx.DeleteNamespace(ctx, o)
}
