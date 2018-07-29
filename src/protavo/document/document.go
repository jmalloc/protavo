package document

import (
	"time"

	"github.com/golang/protobuf/proto"
)

// Document is a document stored in a DB.
type Document struct {
	// ID is the document's unique identifier. It can be any non-empty string.
	ID string

	// Revision is the version of the document as represented by this value.
	// When modifying a document, it must be equal to the revision of the document
	// that is currently persisted, otherwise an OptimisticLockError occurs.
	Revision uint64

	// Keys is the set of indexing keys applied to the document. Keys are used
	// to quickly find a document or set of documents based on identifiers
	// other than the document ID.
	Keys KeyMap

	// Headers is an arbitrary set of key/value pairs that is persisted along
	// with the document content.
	Headers map[string]string

	// CreatedAt is the time at which the document was created. The value is
	// set automatically when the document is saved for the first time.
	CreatedAt time.Time

	// UpdatedAt is the time at which the document was last modified. The value
	// is set automatically when the document is saved.
	UpdatedAt time.Time

	// Content is the application-defined document content.
	Content proto.Message
}

// UniqueKeys returns the document's unique indexing keys.
func (d *Document) UniqueKeys() []string {
	return d.KeysByType(UniqueKey)
}

// SharedKeys returns the document's shared indexing keys.
func (d *Document) SharedKeys() []string {
	return d.KeysByType(SharedKey)
}

// KeysByType returns the subset of the document's keys that are of the given
// type.
func (d *Document) KeysByType(t KeyType) []string {
	keys := make([]string, 0, len(d.Keys))

	for k, kt := range d.Keys {
		if kt == t {
			keys = append(keys, k)
		}
	}

	return keys
}

// Equal returns true if d and doc are equal.
func (d *Document) Equal(doc *Document) bool {
	if d.ID != d.ID {
		return false
	}

	if d.Revision != d.Revision {
		return false
	}

	if !d.CreatedAt.Equal(doc.CreatedAt) {
		return false
	}

	if !d.UpdatedAt.Equal(doc.UpdatedAt) {
		return false
	}

	if len(d.Keys) != len(doc.Keys) {
		return false
	}

	if len(d.Headers) != len(doc.Headers) {
		return false
	}

	for k, v := range d.Keys {
		x, ok := doc.Keys[k]
		if !ok || x != v {
			return false
		}
	}

	for k, v := range d.Headers {
		x, ok := doc.Headers[k]
		if !ok || x != v {
			return false
		}
	}

	return proto.Equal(d.Content, doc.Content)
}

// Validate panics if the document is not valid.
func (d *Document) Validate() {
	if d == nil {
		panic("document must not be nil")
	}

	if d.ID == "" {
		panic("document ID must not be empty")
	}

	if d.Content == nil {
		panic("document content must not be nil")
	}
}
