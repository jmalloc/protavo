package database

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

// GetContent gets the content of the document with the given ID.
func (s *Store) GetContent(id string) (*Content, error) {
	buf := s.Content.Get([]byte(id))
	if buf == nil {
		return nil, fmt.Errorf(
			"data integrity error: content for '%s' is missing",
			id,
		)
	}

	var c Content
	if err := proto.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// PutContent saves the content for a document.
func (s *Store) PutContent(id string, c *Content) error {
	buf, err := proto.Marshal(c)
	if err != nil {
		return err
	}

	return s.Content.Put(
		[]byte(id),
		buf,
	)
}

// DeleteContent deletes the content for the document with the given ID.
func (s *Store) DeleteContent(id string) error {
	return s.Content.Delete(
		[]byte(id),
	)
}
