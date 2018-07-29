package document

import "github.com/golang/protobuf/proto"

// StringContent is a convenience function that returns a Protocol Buffers
// message that contains a single string.
//
// It be used to easily store plain string content inside a document. The string
// can be retreived with GetStringContent().
func StringContent(s string) proto.Message {
	return &StringContentType{Value: s}
}

// GetStringContent returns the string value of a message created with
// StringContent().
func GetStringContent(m proto.Message) string {
	return m.(*StringContentType).Value
}
