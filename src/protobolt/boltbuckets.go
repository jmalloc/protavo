package protobolt

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
)

var (
	rootBucket         = []byte("protobolt")
	metaDataBucketName = []byte("meta")
	contentBucketName  = []byte("content")
	keysBucketName     = []byte("keys")
)

type buckets struct {
	MetaData *bolt.Bucket
	Content  *bolt.Bucket
	Keys     *bolt.Bucket
}

// getBuckets returns all of the buckets used for storage of the given
// namespace.
func getBuckets(ns string, tx *bolt.Tx) (buckets, bool, error) {
	var b buckets

	root := tx.Bucket(rootBucket)
	if root == nil {
		return b, false, nil
	}

	parent := root.Bucket([]byte(ns))
	if parent == nil {
		return b, false, nil
	}

	b.MetaData = parent.Bucket(metaDataBucketName)
	if b.MetaData == nil {
		return b, false, fmt.Errorf(
			"data integrity error: missing '%s' bucket within '%s' namespace",
			metaDataBucketName,
			ns,
		)
	}

	b.Content = parent.Bucket(contentBucketName)
	if b.Content == nil {
		return b, false, fmt.Errorf(
			"data integrity error: missing '%s' bucket within '%s' namespace",
			contentBucketName,
			ns,
		)
	}

	b.Keys = parent.Bucket(keysBucketName)
	if b.Keys == nil {
		return b, false, fmt.Errorf(
			"data integrity error: missing '%s' bucket within '%s' namespace",
			keysBucketName,
			ns,
		)
	}

	return b, true, nil
}

// createBuckets returns all of the buckets used for storage in the given
// namespace, creating them if they don't exist.
func createBuckets(ns string, tx *bolt.Tx) (buckets, error) {
	var b buckets

	root, err := tx.CreateBucketIfNotExists(rootBucket)
	if err != nil {
		return b, err
	}

	parent, err := root.CreateBucketIfNotExists([]byte(ns))
	if err != nil {
		return b, err
	}

	b.MetaData, err = parent.CreateBucketIfNotExists(metaDataBucketName)
	if err != nil {
		return b, err
	}

	b.Content, err = parent.CreateBucketIfNotExists(contentBucketName)
	if err != nil {
		return b, err
	}

	b.Keys, err = parent.CreateBucketIfNotExists(keysBucketName)

	return b, err
}

// deleteBuckets deletes all of the buckets used for storage of the given
// namespace.
func deleteBuckets(ns string, tx *bolt.Tx) error {
	root := tx.Bucket(rootBucket)
	if root == nil {
		return nil
	}

	err := root.DeleteBucket([]byte(ns))

	if err == bolt.ErrBucketNotFound {
		return nil
	}

	return err
}
