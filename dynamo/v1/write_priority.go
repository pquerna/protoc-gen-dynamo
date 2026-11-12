package v1

import "google.golang.org/protobuf/proto"

// Resolved returns p with UNSPECIFIED mapped to NORMAL. This is the single
// definition of the default, shared by the generator and pkg/dynamoattr.
func (p WritePriority) Resolved() WritePriority {
	if p == WritePriority_WRITE_PRIORITY_UNSPECIFIED {
		return WritePriority_WRITE_PRIORITY_NORMAL
	}
	return p
}

// AttrValue returns the stored string for p, read from its
// (dynamo.v1.attr_value) option after resolving UNSPECIFIED to NORMAL. ok is
// false if p is not a known value or has no attr_value.
func (p WritePriority) AttrValue() (string, bool) {
	p = p.Resolved()
	ev := p.Descriptor().Values().ByNumber(p.Number())
	if ev == nil {
		return "", false
	}
	s, _ := proto.GetExtension(ev.Options(), E_AttrValue).(string)
	return s, s != ""
}
