package protobolt

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
)

var (
	metaDataBucketName = []byte("protobolt-meta")
	contentBucketName  = []byte("protobolt-content")
	keysBucketName     = []byte("protobolt-keys")
)

type buckets struct {
	MetaData *bolt.Bucket
	Content  *bolt.Bucket
	Keys     *bolt.Bucket
}

// getBuckets returns all of the buckets used for storage.
func getBuckets(tx *bolt.Tx) (buckets, bool, error) {
	var b buckets

	b.MetaData = tx.Bucket(metaDataBucketName)
	if b.MetaData == nil {
		return b, false, nil
	}

	b.Content = tx.Bucket(contentBucketName)
	if b.Content == nil {
		return b, false, fmt.Errorf(
			"data integrity error: '%s' bucket exists, but '%s' bucket does not",
			metaDataBucketName,
			contentBucketName,
		)
	}

	b.Keys = tx.Bucket(keysBucketName)
	if b.Keys == nil {
		return b, false, fmt.Errorf(
			"data integrity error: '%s' bucket exists, but '%s' bucket does not",
			metaDataBucketName,
			keysBucketName,
		)
	}

	return b, true, nil
}

// createBuckets returns all of the buckets used for storage, creating them
// if they don't exist.
func createBuckets(tx *bolt.Tx) (buckets, error) {
	var (
		b   buckets
		err error
	)

	b.MetaData, err = tx.CreateBucketIfNotExists(metaDataBucketName)
	if err != nil {
		return b, err
	}

	b.Content, err = tx.CreateBucketIfNotExists(contentBucketName)
	if err != nil {
		return b, err
	}

	b.Keys, err = tx.CreateBucketIfNotExists(keysBucketName)

	return b, err
}
