package protobolt

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/protobolt/src/protobolt/internal/types"
)

// getContent gets the content of the document with the given ID.
func getContent(b buckets, docID string) (proto.Message, error) {
	buf := b.Content.Get([]byte(docID))
	if buf == nil {
		return nil, fmt.Errorf(
			"data integrity error: content for '%s' is missing",
			docID,
		)
	}

	var env types.ContentEnvelope
	if err := proto.Unmarshal(buf, &env); err != nil {
		return nil, err
	}

	var any ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(env.Content, &any); err != nil {
		return nil, err
	}

	return any.Message, nil
}

// putContent saves the content for a document.
func putContent(b buckets, docID string, c proto.Message) error {
	any, err := ptypes.MarshalAny(c)
	if err != nil {
		return err
	}

	env := &types.ContentEnvelope{
		Content: any,
	}

	buf, err := proto.Marshal(env)
	if err != nil {
		return err
	}

	return b.Content.Put([]byte(docID), buf)
}
