package database

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
)

const (
	// SharedKeyType is a plain uint32 representation of the document.SharedKey
	// constant.
	SharedKeyType = uint32(document.SharedKey)

	// UniqueKeyType is a plain uint32 representation of the document.UniqueKey
	// constant.
	UniqueKeyType = uint32(document.UniqueKey)
)

// UpdateKeys updates the unique keys for a specific document.
func (s *Store) UpdateKeys(
	id string,
	before, after map[string]uint32,
) error {
	// remove the document from any keys that are not present after the update
	for key := range before {
		if _, ok := after[key]; ok {
			continue
		}

		k, err := s.GetKey(key)
		if err != nil {
			return err
		}

		delete(k.Documents, id)

		if err := s.putKey(key, k); err != nil {
			return err
		}
	}

	// add the document to any keys that are present after the update
	for key, afterType := range after {
		// optimization: skip over any keys that haven't changed at all
		if beforeType, ok := before[key]; ok {
			if beforeType == afterType {
				continue
			}
		}

		k, err := s.GetKey(key)
		if err != nil {
			return err
		}

		// count how many other docs are in this key
		otherDocs := len(k.Documents)
		if k.Documents[id] {
			otherDocs--
		}

		if otherDocs == 0 {
			// there are no other documents, update the key however neceessary
			k.Type = afterType
			k.Documents = map[string]bool{id: true}
		} else if k.Type == SharedKeyType && afterType == SharedKeyType {
			// there are other documents, but the key is shared AND the document wants a
			// shared key, so simply add the document to the key
			k.Documents[id] = true
		} else {
			// otherwise, either the key is already unique, or the document is requesting
			// a unique key, hence there is a conflict
			for conflictID := range k.Documents {
				return &protavo.DuplicateKeyError{
					DocumentID:            id,
					ConflictingDocumentID: conflictID,
					UniqueKey:             key,
				}
			}

			panic("impossible condition: len(k.Documents) > 0 but iterating it produces no values")
		}

		if err := s.putKey(key, k); err != nil {
			return err
		}
	}

	return nil
}

// GetKey loads a key by its name.
// If the key does not exist, a pointer to a zero-value Key is returned.
func (s *Store) GetKey(key string) (*Key, error) {
	var k Key

	if buf := s.Keys.Get([]byte(key)); buf != nil {
		if err := proto.Unmarshal(buf, &k); err != nil {
			return nil, err
		}
	}

	return &k, nil
}

// putKey saves a key, deleting it if it contains no documents.
func (s *Store) putKey(key string, k *Key) error {
	if len(k.Documents) == 0 {
		return s.Keys.Delete([]byte(key))
	}

	buf, err := proto.Marshal(k)
	if err != nil {
		return err
	}

	return s.Keys.Put([]byte(key), buf)
}
