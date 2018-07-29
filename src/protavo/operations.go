package protavo

import (
	"github.com/jmalloc/protavo/src/protavo/driver"
	"github.com/jmalloc/protavo/src/protavo/filter"
)

// FetchAll invokes fn for every document in the store.
func FetchAll(fn driver.FetchFunc) *driver.Fetch {
	return &driver.Fetch{
		Each: fn,
	}
}

// FetchWhere invokes fn for each document that matches the all of the given
// filter conditions.
func FetchWhere(fn driver.FetchFunc, f ...filter.Condition) *driver.Fetch {
	return &driver.Fetch{
		Each:   fn,
		Filter: filter.New(f),
	}
}

// Save saves a document to the store, creating it if it does not already exist.
//
// The revision of the document must match the current revision of that document
// in the store; otherwise, an OptimisticLockError is returned.
func Save(doc *Document) *driver.Save {
	return &driver.Save{
		Document: doc,
	}
}

// ForceSave saves a document to the store without checking the current
// revision.
func ForceSave(doc *Document) *driver.Save {
	return &driver.Save{
		Document: doc,
		Force:    true,
	}
}

// Delete removes a document from the store.
//
// The revision of the document must match the current revision of that document
// in the store; otherwise, an OptimistcLockError is returned.
//
// It is not an error to delete a document that does not exist, provided that
// the given revision is zero.
func Delete(doc *Document) *driver.Delete {
	return &driver.Delete{
		Document: doc,
	}
}

// DeleteWhere removes documents from the store without checking the
// current document revisions.
//
// It is not an error to delete a document that does not exist.
func DeleteWhere(f ...filter.Condition) *driver.DeleteWhere {
	return &driver.DeleteWhere{
		Filter: filter.New(f),
	}
}
