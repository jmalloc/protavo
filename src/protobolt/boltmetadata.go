package protobolt

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/protobolt/src/protobolt/internal/types"
)

// checkVersion returns an optimistic lock error if given and actual do not have
// the same version.
func checkVersion(action string, before, after *types.MetaData) error {
	var ver uint64
	if before != nil {
		ver = before.Version
	}

	if after.Version == ver {
		return nil
	}

	return &OptimisticLockError{
		DocumentID:    after.Id,
		GivenVersion:  after.Version,
		ActualVersion: ver,
		Action:        action,
	}
}

// tryGetMetaData unmarshals the meta-data in b with key k.
func tryGetMetaData(b buckets, docID string) (*types.MetaData, bool, error) {
	buf := b.MetaData.Get([]byte(docID))
	if buf == nil {
		return nil, false, nil
	}

	md, err := unmarshalMetaData(buf)
	return md, true, err
}

// getMetaData unmarshals the meta-data in b with key k.
// it returns an error if the meta-data does not exist.
func getMetaData(b buckets, docID string) (*types.MetaData, error) {
	md, ok, err := tryGetMetaData(b, docID)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf(
			"data integrity error: meta-data for '%s' is missing",
			docID,
		)
	}

	return md, nil
}

// putMetaData saves content to the given bucket.
func putMetaData(b buckets, md *types.MetaData) error {
	buf, err := proto.Marshal(md)
	if err != nil {
		return err
	}

	return b.MetaData.Put([]byte(md.Id), buf)
}

// unmarshalMetaData unmarshals document meta-data from its binary representation.
func unmarshalMetaData(buf []byte) (*types.MetaData, error) {
	var md types.MetaData

	if err := proto.Unmarshal(buf, &md); err != nil {
		return nil, err
	}

	return &md, nil
}
