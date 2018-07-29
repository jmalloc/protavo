package driver

import (
	"context"

	"github.com/jmalloc/protavo/src/protavo/document"
)

// Save is a request to save a document.
type Save struct {
	operation

	Document *document.Document
	Force    bool
}

// ExecuteInWriteTx executes this operation within the context of tx.
func (o *Save) ExecuteInWriteTx(ctx context.Context, tx WriteTx) {
	tx.Save(ctx, o)
}
