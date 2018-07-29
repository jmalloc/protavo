package protavobolt

import (
	bolt "github.com/coreos/bbolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
	"github.com/jmalloc/protavo/src/protavobolt/internal/database"
)

// executeSave creates or updates a set of documents.
func executeSave(
	tx *bolt.Tx,
	ns string,
	doc *document.Document,
	force bool,
) error {
	s, err := database.CreateStore(tx, ns)
	if err != nil {
		return err
	}

	rec, exists, err := s.TryGetRecord(doc.ID)
	if err != nil {
		return err
	}

	if !force && doc.Revision != rec.GetRevision() {
		return &protavo.OptimisticLockError{
			DocumentID: doc.ID,
			GivenRev:   doc.Revision,
			ActualRev:  rec.GetRevision(),
			Operation:  "save",
		}
	}

	var new *database.Record
	if exists {
		new, err = updateRecord(s, doc, rec)
	} else {
		new, err = createRecord(s, doc)
	}
	if err != nil {
		return err
	}

	c, err := marshalContent(doc)
	if err != nil {
		return err
	}

	if err := s.PutContent(doc.ID, c); err != nil {
		return err
	}

	return unmarshalRecordManagedFields(new, doc)
}

// createRecord creates a new document record.
func createRecord(
	s *database.Store,
	doc *document.Document,
) (*database.Record, error) {
	now := ptypes.TimestampNow()
	new := &database.Record{
		Revision:  1,
		Keys:      marshalKeys(doc.Keys),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.PutRecord(doc.ID, new); err != nil {
		return nil, err
	}

	if err := s.UpdateKeys(doc.ID, nil, new.Keys); err != nil {
		return nil, err
	}

	return new, nil
}

// updateRecord updates an existing document record.
func updateRecord(
	s *database.Store,
	doc *document.Document,
	rec *database.Record,
) (*database.Record, error) {
	new := proto.Clone(rec).(*database.Record)
	new.Revision++
	new.Keys = marshalKeys(doc.Keys)
	new.UpdatedAt = ptypes.TimestampNow()

	if err := s.PutRecord(doc.ID, new); err != nil {
		return nil, err
	}

	if err := s.UpdateKeys(doc.ID, rec.Keys, new.Keys); err != nil {
		return nil, err
	}

	return new, nil
}
