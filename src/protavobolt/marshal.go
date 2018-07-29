package protavobolt

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// newDocument constructs a new document from a record and content.
func newDocument(
	id string,
	rec *database.Record,
	c *database.Content,
) (*document.Document, error) {
	doc := &document.Document{
		ID:   id,
		Keys: unmarshalKeys(rec.Keys),
	}

	if err := unmarshalContent(c, doc); err != nil {
		return nil, nil
	}

	return doc, unmarshalRecordManagedFields(rec, doc)
}

// unmarshalRecordManagedFields synchronizes a document with data from a record.
// It only updates the fields that are populated by the implementation, as
// opposed to the user.
func unmarshalRecordManagedFields(rec *database.Record, doc *document.Document) error {
	createdAt, err := ptypes.Timestamp(rec.CreatedAt)
	if err != nil {
		return err
	}

	updatedAt, err := ptypes.Timestamp(rec.UpdatedAt)
	if err != nil {
		return err
	}

	doc.Revision = rec.Revision
	doc.CreatedAt = createdAt
	doc.UpdatedAt = updatedAt

	return nil
}

// marshalContent converts content from the public API to database format.
func marshalContent(doc *document.Document) (*database.Content, error) {
	c := &database.Content{
		Headers: doc.Headers,
	}

	var err error
	c.Content, err = ptypes.MarshalAny(doc.Content)
	return c, err
}

// unmarshalContent converts content from the databsae to public API format.
func unmarshalContent(c *database.Content, doc *document.Document) error {
	var x ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(c.Content, &x); err != nil {
		return err
	}

	doc.Headers = c.Headers
	doc.Content = x.Message

	return nil
}

// marshalKeys converts a key map from the public API to database format.
func marshalKeys(keys map[string]document.KeyType) map[string]uint32 {
	r := make(map[string]uint32, len(keys))

	for k, v := range keys {
		r[k] = uint32(v)
	}

	return r
}

// unmarshalKeys converts a key map from the database to public API format.
func unmarshalKeys(keys map[string]uint32) map[string]document.KeyType {
	r := make(map[string]document.KeyType, len(keys))

	for k, v := range keys {
		r[k] = document.KeyType(v)
	}

	return r
}
