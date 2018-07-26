package protobolt

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/protobolt/src/protobolt/internal/types"
)

// Document is a document stored in a DB.
type Document struct {
	md *types.MetaData
	c  proto.Message
}

// NewDocument returns a new document with the given content.
func NewDocument(id string, content proto.Message) *Document {
	d := &Document{
		&types.MetaData{
			Id: id,
		},
		content,
	}

	d.validate()

	return d
}

// ID returns the unique identifier of the document.
func (d *Document) ID() string {
	return d.md.Id
}

// Headers returns a copy of the document's headers.
func (d *Document) Headers() map[string]string {
	h := map[string]string{}

	for k, v := range d.md.Headers {
		h[k] = v
	}

	return h
}

// Version returns the version of the document.
func (d *Document) Version() uint64 {
	return d.md.Version
}

// CreatedAt returns the time at which the document was first saved.
// It returns false if the document has not yet been saved.
func (d *Document) CreatedAt() (time.Time, bool) {
	if d.md.Version == 0 {
		return time.Time{}, false
	}

	t, err := ptypes.Timestamp(d.md.CreatedAt)
	if err != nil {
		return time.Time{}, false
	}

	return t, true
}

// UpdatedAt returns the time at which the document was last saved.
// It returns false if the document has not yet been saved.
func (d *Document) UpdatedAt() (time.Time, bool) {
	if d.md.Version == 0 {
		return time.Time{}, false
	}

	t, err := ptypes.Timestamp(d.md.UpdatedAt)
	if err != nil {
		return time.Time{}, false
	}

	return t, true
}

// Content returns a copy of the document's meta-data.
func (d *Document) Content() proto.Message {
	return proto.Clone(d.c)
}

// Equal returns true if d and doc are equal.
func (d *Document) Equal(doc *Document) bool {
	return proto.Equal(d.md, doc.md) && proto.Equal(d.c, doc.c)
}

// WithContent returns a copy of this document containing the given content.
func (d *Document) WithContent(content proto.Message) *Document {
	if proto.Equal(content, d.c) {
		return d
	}

	return &Document{
		d.md,
		proto.Clone(d.c),
	}
}

// WithHeader returns a copy of this document which includes a new header.
func (d *Document) WithHeader(k, v string) *Document {
	md := d.cloneMetaData()

	if md.Headers == nil {
		md.Headers = map[string]string{}
	}

	md.Headers[k] = v

	return &Document{md, d.c}
}

// WithHeaders returns a copy of this document with the given headers merged in
// to the existing set of headers.
func (d *Document) WithHeaders(headers map[string]string) *Document {
	if len(headers) == 0 {
		return d
	}

	md := d.cloneMetaData()

	if md.Headers == nil {
		md.Headers = headers
	} else {
		for k, v := range headers {
			md.Headers[k] = v
		}
	}

	return &Document{md, d.c}
}

// WithoutHeaders returns a copy of this document without the given headers.
func (d *Document) WithoutHeaders(headers ...string) *Document {
	if len(headers) == 0 || len(d.md.Headers) == 0 {
		return d
	}

	md := d.cloneMetaData()

	for _, k := range headers {
		delete(md.Headers, k)
	}

	return &Document{md, d.c}
}

// withKeys returns a copy of this document with the given keys.
func (d *Document) withKeys(t types.KeyType, keys ...string) *Document {
	if len(keys) == 0 {
		return d
	}

	md := d.cloneMetaData()

	if md.Keys == nil {
		md.Keys = map[string]types.KeyType{}
	}

	for _, k := range keys {
		md.Keys[k] = t
	}

	return &Document{md, d.c}
}

// WithUniqueKeys returns a copy of this document with the given unique keys.
func (d *Document) WithUniqueKeys(keys ...string) *Document {
	return d.withKeys(types.KeyType_UNIQUE, keys...)
}

// WithSharedKeys returns a copy of this document with the given shared keys.
func (d *Document) WithSharedKeys(keys ...string) *Document {
	return d.withKeys(types.KeyType_SHARED, keys...)
}

// WithoutKeys returns a copy of this document without the given keys.
func (d *Document) WithoutKeys(keys ...string) *Document {
	if len(keys) == 0 || len(d.md.Keys) == 0 {
		return d
	}

	md := d.cloneMetaData()

	for _, k := range keys {
		delete(md.Keys, k)
	}

	return &Document{md, d.c}
}

func (d *Document) cloneMetaData() *types.MetaData {
	return proto.Clone(d.md).(*types.MetaData)
}

func (d *Document) validate() {
	if d == nil {
		panic("document must not be nil")
	}

	if d.md.Id == "" {
		panic("document ID must not be empty")
	}

	if d.c == nil {
		panic("document content must not be empty")
	}
}
