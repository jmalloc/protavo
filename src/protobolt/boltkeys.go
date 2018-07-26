package protobolt

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/protobolt/src/protobolt/internal/types"
)

// matchesFilter return true if the document with the given meta data has
// all of the given keys.
func matchesFilter(md *types.MetaData, filter []string) bool {
	for _, key := range filter {
		if _, ok := md.Keys[key]; !ok {
			return false
		}
	}

	return true
}

// updateKeys updates the unique keys for a specific document.
func updateKeys(b buckets, docID string, before, after *types.MetaData) error {
	beforeKeys := before.GetKeys()
	afterKeys := after.GetKeys()

	// remove the document from any keys that are not present after the update
	for key := range beforeKeys {
		if _, ok := afterKeys[key]; ok {
			continue
		}

		kd, err := getKeyData(b, key)
		if err != nil {
			return err
		}

		delete(kd.Documents, docID)

		if err := putKeyData(b, key, kd); err != nil {
			return err
		}
	}

	// add the document to any keys that are present after the update
	for key, afterType := range afterKeys {
		// optimization: skip over any keys that haven't changed at all
		if beforeType, ok := beforeKeys[key]; ok {
			if beforeType == afterType {
				continue
			}
		}

		kd, err := getKeyData(b, key)
		if err != nil {
			return err
		}

		// count how many other docs are in this key
		otherDocs := len(kd.Documents)
		if kd.Documents[docID] {
			otherDocs--
		}

		if otherDocs == 0 {
			// there are no other documents, update the key however neceessary
			kd.Documents = map[string]bool{docID: true}
			kd.Type = afterType
		} else if kd.Type == types.KeyType_SHARED && afterType == types.KeyType_SHARED {
			// there are other documents, but the key is shared AND the document wants a
			// shared key, so simply add the document to the key
			kd.Documents[docID] = true
		} else {
			// otherwise, either the key is already unique, or the document is requesting
			// a unique key, hence there is a conflict
			for conflictID := range kd.Documents {
				return &DuplicateKeyError{
					DocumentID:            docID,
					ConflictingDocumentID: conflictID,
					UniqueKey:             key,
				}
			}

			panic("impossible condition: len(kd.Documents) > 0 but iterating it produces no values")
		}

		if err := putKeyData(b, key, kd); err != nil {
			return err
		}
	}

	return nil
}

// getKeyData loads the KeyData for the given key.
// If the key does not exist, a pointer to a zero-value KeyData is returned.
func getKeyData(b buckets, key string) (*types.KeyData, error) {
	var kd types.KeyData

	if buf := b.Keys.Get([]byte(key)); buf != nil {
		if err := proto.Unmarshal(buf, &kd); err != nil {
			return nil, err
		}
	}

	return &kd, nil
}

// putKeyData saves the KeyData for the given key, deleting it if it contains
// no documents.
func putKeyData(b buckets, key string, kd *types.KeyData) error {
	if len(kd.Documents) == 0 {
		return b.Keys.Delete([]byte(key))
	}

	buf, err := proto.Marshal(kd)
	if err != nil {
		return err
	}

	return b.Keys.Put([]byte(key), buf)
}
