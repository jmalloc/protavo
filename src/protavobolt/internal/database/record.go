package database

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

// TryGetRecord loads the record for the document with the given ID.
//
// The presence of the record is the authoratative indicator of a document's
// existence. It returns false if the document is not in the store.
func (s *Store) TryGetRecord(id string) (*Record, bool, error) {
	buf := s.Records.Get([]byte(id))
	if buf == nil {
		return nil, false, nil
	}

	rec, err := UnmarshalRecord(buf)
	return rec, true, err
}

// GetRecord loads the record for the document with the gien ID.
//
// It is intended for use when the document is known to exist.
// It returns an error if the document is not in the store.
func (s *Store) GetRecord(id string) (*Record, error) {
	rec, ok, err := s.TryGetRecord(id)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf(
			"data integrity error: record for '%s' is missing",
			id,
		)
	}

	return rec, nil
}

// PutRecord creates or updates a document record.
func (s *Store) PutRecord(id string, rec *Record) error {
	buf, err := proto.Marshal(rec)
	if err != nil {
		return err
	}

	return s.Records.Put(
		[]byte(id),
		buf,
	)
}

// DeleteRecord deletes a document record.
func (s *Store) DeleteRecord(id string) error {
	return s.Records.Delete(
		[]byte(id),
	)
}

// UnmarshalRecord unmarshals document record.
func UnmarshalRecord(buf []byte) (*Record, error) {
	var rec Record

	if err := proto.Unmarshal(buf, &rec); err != nil {
		return nil, err
	}

	return &rec, nil
}
