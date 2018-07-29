package protavo

import (
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
)

// FetchAll returns an operation that calls fn once for every document.
//
// It stops iterating if fn returns false or a non-nil error.
func FetchAll(fn driver.FetchFunc) *driver.Fetch {
	return &driver.Fetch{
		Each: fn,
	}
}

// FetchWhere returns an operation that calls fn once for each document that
// matches the given filter conditions.
//
// It stops iterating if fn returns false or a non-nil error.
func FetchWhere(fn driver.FetchFunc, f ...filter.Condition) *driver.Fetch {
	return &driver.Fetch{
		Each:   fn,
		Filter: filter.New(f),
	}
}

// Save atomically creates or update a documents.
//
// The `Revision` field of the document must be equal to the revision of that
// document as currently persisted; otherwise, an OptimisticLockError is
// returned.
//
// New documents must have a `Revision` of `0`.
//
// doc is updated with its new revision and timestamp.
func Save(doc *Document) *driver.Save {
	return &driver.Save{
		Document: doc,
	}
}

// ForceSave creates or updates a documents without checking the current
// revisions.
//
// doc is updated with its new revision and timestamp.
func ForceSave(doc *Document) *driver.Save {
	return &driver.Save{
		Document: doc,
		Force:    true,
	}
}

// Delete removes a document.
//
// The `Revision` field of the document must be equal to the revision of that
// document as currently persisted; otherwise, an OptimisticLockError is
// returned.
//
// It is not an error to delete a non-existent document, provided the given
// `Revision` is `0`.
func Delete(doc *Document) *driver.Delete {
	return &driver.Delete{
		Document: doc,
	}
}

// DeleteWhere atomically removes the documents that match the given filter
// conditions without checking the current revisions.
//
// If fn is non-nil, it is invoked for each of the deleted documents.
func DeleteWhere(
	fn driver.DeleteWhereFunc,
	f ...filter.Condition,
) *driver.DeleteWhere {
	return &driver.DeleteWhere{
		Filter: filter.New(f),
	}
}
