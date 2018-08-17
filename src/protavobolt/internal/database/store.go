package database

import (
	"bytes"
	"fmt"

	bolt "github.com/coreos/bbolt"
)

var (
	rootBucket    = []byte("protavo")
	recordsBucket = []byte("records")
	contentBucket = []byte("content")
	keysBucket    = []byte("keys")
)

// Store is the data store for a single namespace.
type Store struct {
	Records *bolt.Bucket
	Content *bolt.Bucket
	Keys    *bolt.Bucket
}

// OpenStore returns the store for the given namespace.
//
// It returns false if the store does not exist.
func OpenStore(tx *bolt.Tx, ns string) (*Store, bool, error) {
	parent := tx.Bucket(rootBucket)
	if parent == nil {
		return nil, false, nil
	}

	if ns != "" {
		for _, p := range splitNamespace(ns) {
			parent = parent.Bucket(p)
			if parent == nil {
				return nil, false, nil
			}
		}
	}

	s := &Store{}

	s.Records = parent.Bucket(recordsBucket)
	if s.Records == nil {
		return nil, false, fmt.Errorf(
			"data integrity error: missing '%s' bucket within '%s' namespace",
			recordsBucket,
			ns,
		)
	}

	s.Content = parent.Bucket(contentBucket)
	if s.Content == nil {
		return nil, false, fmt.Errorf(
			"data integrity error: missing '%s' bucket within '%s' namespace",
			contentBucket,
			ns,
		)
	}

	s.Keys = parent.Bucket(keysBucket)
	if s.Keys == nil {
		return nil, false, fmt.Errorf(
			"data integrity error: missing '%s' bucket within '%s' namespace",
			keysBucket,
			ns,
		)
	}

	return s, true, nil
}

// CreateStore returns the store for a single namespace, creating it if it does
// not exist.
func CreateStore(tx *bolt.Tx, ns string) (*Store, error) {
	parent, err := tx.CreateBucketIfNotExists(rootBucket)
	if err != nil {
		return nil, err
	}

	if ns != "" {
		for _, p := range splitNamespace(ns) {
			parent, err = parent.CreateBucketIfNotExists(p)
			if err != nil {
				return nil, err
			}
		}
	}

	s := &Store{}

	s.Records, err = parent.CreateBucketIfNotExists(recordsBucket)
	if err != nil {
		return nil, err
	}

	s.Content, err = parent.CreateBucketIfNotExists(contentBucket)
	if err != nil {
		return nil, err
	}

	s.Keys, err = parent.CreateBucketIfNotExists(keysBucket)

	return s, err
}

// DeleteStore deletes the store for a given namespace.
// It is not an error to delete a non-existent store.
func DeleteStore(tx *bolt.Tx, ns string) error {
	var parent interface {
		Bucket([]byte) *bolt.Bucket
		DeleteBucket([]byte) error
	} = tx

	name := rootBucket

	if ns != "" {
		parent = parent.Bucket(rootBucket)
		if parent == nil {
			return nil
		}
		parts := splitNamespace(ns)
		last := len(parts) - 1
		name = []byte(parts[last])

		for _, p := range parts[:last] {
			parent = parent.Bucket([]byte(p))
			if parent == nil {
				return nil
			}
		}
	}

	err := parent.DeleteBucket(name)

	if err == bolt.ErrBucketNotFound {
		return nil
	}

	return err
}

func splitNamespace(ns string) [][]byte {
	return bytes.Split([]byte(ns), []byte("."))
}
