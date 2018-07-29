package protavo

import (
	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
)

// FetchAll returns an operation that calls fn once for every document.
//
// It stops iterating if fn returns false or a non-nil error.
//
// The returned operation can be executed atomically with other operations using
// DB.Read() or DB.Write().
func FetchAll(fn driver.FetchFunc) *driver.Fetch {
	return &driver.Fetch{
		Each: fn,
	}
}

// FetchWhere returns an operation that calls fn once for each document that
// matches the given filter conditions.
//
// It stops iterating if fn returns false or a non-nil error.
//
// The returned operation can be executed atomically with other operations using
// DB.Read() or DB.Write().
func FetchWhere(fn driver.FetchFunc, f ...filter.Condition) *driver.Fetch {
	return &driver.Fetch{
		Each:   fn,
		Filter: filter.New(f),
	}
}

// Save returns an operation that creates or update a document.
//
// The Revision field of the document must be equal to the revision of that
// document as currently persisted; otherwise, an OptimisticLockError is
// returned.
//
// New documents must have a revision of 0.
//
// doc is updated with its new revision and timestamp.
//
// The returned operation can be executed atomically with other operations using
// DB.Write().
func Save(doc *document.Document) *driver.Save {
	return &driver.Save{
		Document: doc,
	}
}

// ForceSave returns an operation that creates or updates a document without
// checking the current revision.
//
// doc is updated with its new revision and timestamp.
//
// The returned operation can be executed atomically with other operations using
// DB.Write().
func ForceSave(doc *document.Document) *driver.Save {
	return &driver.Save{
		Document: doc,
		Force:    true,
	}
}

// Delete returns an operation that removes a document.
//
// The Revision field of the document must be equal to the revision of that
// document as currently persisted; otherwise, an OptimisticLockError is
// returned.
//
// It is not an error to delete a non-existent document, provided the given
// revision is 0.
//
// The returned operation can be executed atomically with other operations using
// DB.Write().
func Delete(doc *document.Document) *driver.Delete {
	return &driver.Delete{
		Document: doc,
	}
}

// DeleteWhere returns an operation that atomically removes the documents that
// match the given filter conditions without checking the current revisions.
//
// If fn is non-nil, it is invoked for each of the deleted documents.
//
// The returned operation can be executed atomically with other operations using
// DB.Write().
func DeleteWhere(
	fn driver.DeleteWhereFunc,
	f ...filter.Condition,
) *driver.DeleteWhere {
	return &driver.DeleteWhere{
		Filter: filter.New(f),
	}
}
