// Package dynamoattr reads and writes the item-level attributes that
// protoc-gen-dynamo emits outside the marshalled message.
package dynamoattr

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	dynamov1 "github.com/pquerna/protoc-gen-dynamo/dynamo/v1"
)

// WritePriorityAttr is the DynamoDB attribute name holding the write priority.
//
// The priority describes how the change carrying it should be processed, not
// the object: it is stamped on the item only so change-stream consumers can
// see it on the stream record.
const WritePriorityAttr = "wp"

var attrToPriority = map[string]dynamov1.WritePriority{}

// The reverse map is built from the attr_value enum-value option, the same
// source the generator reads, so the two cannot disagree.
func init() {
	values := dynamov1.WritePriority(0).Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		p := dynamov1.WritePriority(values.Get(i).Number())
		if p == dynamov1.WritePriority_WRITE_PRIORITY_UNSPECIFIED {
			continue
		}
		if s, ok := p.AttrValue(); ok {
			attrToPriority[s] = p
		}
	}
}

// WritePriorityAttrValue returns the stored string for p. UNSPECIFIED maps to
// the NORMAL value. A value unknown to this build (e.g. an out-of-range cast)
// also maps to NORMAL.
func WritePriorityAttrValue(p dynamov1.WritePriority) string {
	if s, ok := p.AttrValue(); ok {
		return s
	}
	s, _ := dynamov1.WritePriority_WRITE_PRIORITY_NORMAL.AttrValue()
	return s
}

// ParseWritePriority maps a stored string back to the enum. ok is false for an
// empty or unrecognised value (e.g. one written by a newer producer); callers
// must handle that rather than fail.
func ParseWritePriority(s string) (dynamov1.WritePriority, bool) {
	p, ok := attrToPriority[s]
	return p, ok
}

// GetWritePriority reads the attribute from a marshalled item. ok is false if
// the attribute is missing (e.g. an item written before it existed), is not a
// string, or holds an unrecognised value.
func GetWritePriority(item map[string]types.AttributeValue) (dynamov1.WritePriority, bool) {
	av, ok := item[WritePriorityAttr].(*types.AttributeValueMemberS)
	if !ok {
		return dynamov1.WritePriority_WRITE_PRIORITY_UNSPECIFIED, false
	}
	return ParseWritePriority(av.Value)
}
