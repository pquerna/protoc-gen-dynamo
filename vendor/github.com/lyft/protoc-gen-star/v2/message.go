package pgs

import (
	"google.golang.org/protobuf/runtime/protoimpl"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

// Message describes a proto message. Messages can be contained in either
// another Message or File, and may house further Messages and/or Enums. While
// all Fields technically live on the Message, some may be contained within
// OneOf blocks.
type Message interface {
	ParentEntity

	// Descriptor returns the underlying proto descriptor for this message
	Descriptor() *descriptor.DescriptorProto

	// Parent returns either the File or Message that directly contains this
	// Message.
	Parent() ParentEntity

	// Fields returns all fields on the message, including those contained within
	// OneOf blocks.
	Fields() []Field

	// NonOneOfFields returns all fields not contained within OneOf blocks.
	NonOneOfFields() []Field

	// OneOfFields returns only the fields contained within OneOf blocks.
	OneOfFields() []Field

	// SyntheticOneOfFields returns only the fields contained within synthetic OneOf blocks.
	// See: https://github.com/protocolbuffers/protobuf/blob/v3.17.0/docs/field_presence.md
	SyntheticOneOfFields() []Field

	// OneOfs returns the OneOfs contained within this Message.
	OneOfs() []OneOf

	// RealOneOfs returns the OneOfs contained within this Message.
	// This excludes synthetic OneOfs.
	// See: https://github.com/protocolbuffers/protobuf/blob/v3.17.0/docs/field_presence.md
	RealOneOfs() []OneOf

	// Extensions returns all of the Extensions applied to this Message.
	Extensions() []Extension

	// Dependents returns all of the messages where message is directly or
	// transitively used.
	Dependents() []Message

	// Dependencies returns all of the messages that message directly or
	// transitively uses.
	Dependencies() []Message

	// IsMapEntry identifies this message as a MapEntry. If true, this message is
	// not generated as code, and is used exclusively when marshaling a map field
	// to the wire format.
	IsMapEntry() bool

	// IsWellKnown identifies whether or not this Message is a WKT from the
	// `google.protobuf` package. Most official plugins special case these types
	// and they usually need to be handled differently.
	IsWellKnown() bool

	// WellKnownType returns the WellKnownType associated with this field. If
	// IsWellKnown returns false, UnknownWKT is returned.
	WellKnownType() WellKnownType

	setParent(p ParentEntity)
	addField(f Field)
	addExtension(e Extension)
	addOneOf(o OneOf)
	addDependent(message Message)
	getDependents(set map[string]Message)
	addDependency(message Message)
	getDependencies(set map[string]Message)
}

type msg struct {
	desc   *descriptor.DescriptorProto
	parent ParentEntity
	fqn    string

	msgs, preservedMsgs []Message
	enums               []Enum
	exts                []Extension
	defExts             []Extension
	fields              []Field
	oneofs              []OneOf
	maps                []Message
	dependents          []Message
	dependentsCache     map[string]Message
	dependencies        []Message
	dependenciesCache   map[string]Message

	info SourceCodeInfo
}

func (m *msg) Name() Name                              { return Name(m.desc.GetName()) }
func (m *msg) FullyQualifiedName() string              { return m.fqn }
func (m *msg) Syntax() Syntax                          { return m.parent.Syntax() }
func (m *msg) Package() Package                        { return m.parent.Package() }
func (m *msg) File() File                              { return m.parent.File() }
func (m *msg) BuildTarget() bool                       { return m.parent.BuildTarget() }
func (m *msg) SourceCodeInfo() SourceCodeInfo          { return m.info }
func (m *msg) Descriptor() *descriptor.DescriptorProto { return m.desc }
func (m *msg) Parent() ParentEntity                    { return m.parent }
func (m *msg) IsMapEntry() bool                        { return m.desc.GetOptions().GetMapEntry() }
func (m *msg) Enums() []Enum                           { return m.enums }
func (m *msg) Messages() []Message                     { return m.msgs }
func (m *msg) Fields() []Field                         { return m.fields }
func (m *msg) OneOfs() []OneOf                         { return m.oneofs }
func (m *msg) MapEntries() []Message                   { return m.maps }

func (m *msg) WellKnownType() WellKnownType {
	if m.Package().ProtoName() == WellKnownTypePackage {
		return LookupWKT(m.Name())
	}
	return UnknownWKT
}

func (m *msg) IsWellKnown() bool {
	return m.WellKnownType().Valid()
}

func (m *msg) AllEnums() []Enum {
	es := m.Enums()
	for _, m := range m.msgs {
		es = append(es, m.AllEnums()...)
	}
	return es
}

func (m *msg) AllMessages() []Message {
	msgs := m.Messages()
	for _, sm := range m.msgs {
		msgs = append(msgs, sm.AllMessages()...)
	}
	return msgs
}

func (m *msg) NonOneOfFields() (f []Field) {
	for _, fld := range m.fields {
		if !fld.InOneOf() {
			f = append(f, fld)
		}
	}
	return f
}

func (m *msg) OneOfFields() (f []Field) {
	for _, o := range m.oneofs {
		f = append(f, o.Fields()...)
	}

	return f
}

func (m *msg) SyntheticOneOfFields() (f []Field) {
	for _, o := range m.oneofs {
		if o.IsSynthetic() {
			f = append(f, o.Fields()...)
		}
	}

	return f
}

func (m *msg) RealOneOfs() (r []OneOf) {
	for _, o := range m.oneofs {
		if !o.IsSynthetic() {
			r = append(r, o)
		}
	}

	return r
}

func (m *msg) Imports() (i []File) {
	// Mapping for avoiding duplicate entries
	mp := make(map[string]File, len(m.fields))
	for _, f := range m.fields {
		for _, imp := range f.Imports() {
			mp[imp.File().Name().String()] = imp
		}
	}
	for _, f := range mp {
		i = append(i, f)
	}
	return
}

func (m *msg) getDependents(set map[string]Message) {
	m.populateDependentsCache()

	for fqn, d := range m.dependentsCache {
		set[fqn] = d
	}
}

func (m *msg) populateDependentsCache() {
	if m.dependentsCache != nil {
		return
	}

	m.dependentsCache = map[string]Message{}
	for _, dep := range m.dependents {
		m.dependentsCache[dep.FullyQualifiedName()] = dep
		dep.getDependents(m.dependentsCache)
	}
}

func (m *msg) Dependents() []Message {
	m.populateDependentsCache()
	return messageSetToSlice(m.FullyQualifiedName(), m.dependentsCache)
}

func (m *msg) getDependencies(set map[string]Message) {
	m.populateDependenciesCache()

	for fqn, d := range m.dependenciesCache {
		set[fqn] = d
	}
}

func (m *msg) populateDependenciesCache() {
	if m.dependenciesCache != nil {
		return
	}

	m.dependenciesCache = map[string]Message{}
	for _, dep := range m.dependencies {
		m.dependenciesCache[dep.FullyQualifiedName()] = dep
		dep.getDependencies(m.dependenciesCache)
	}
}

func (m *msg) Dependencies() []Message {
	m.populateDependenciesCache()
	return messageSetToSlice(m.FullyQualifiedName(), m.dependenciesCache)
}

func (m *msg) Extension(desc *protoimpl.ExtensionInfo, ext interface{}) (bool, error) {
	return extension(m.desc.GetOptions(), desc, &ext)
}

func (m *msg) Extensions() []Extension {
	return m.exts
}

func (m *msg) DefinedExtensions() []Extension {
	return m.defExts
}

func (m *msg) accept(v Visitor) (err error) {
	if v == nil {
		return nil
	}

	if v, err = v.VisitMessage(m); err != nil || v == nil {
		return
	}

	for _, e := range m.enums {
		if err = e.accept(v); err != nil {
			return
		}
	}

	for _, sm := range m.msgs {
		if err = sm.accept(v); err != nil {
			return
		}
	}

	for _, f := range m.fields {
		if err = f.accept(v); err != nil {
			return
		}
	}

	for _, o := range m.oneofs {
		if err = o.accept(v); err != nil {
			return
		}
	}

	for _, ext := range m.defExts {
		if err = ext.accept(v); err != nil {
			return
		}
	}

	return
}

func (m *msg) addExtension(ext Extension) {
	m.exts = append(m.exts, ext)
}

func (m *msg) addDefExtension(ext Extension) {
	m.defExts = append(m.defExts, ext)
}

func (m *msg) setParent(p ParentEntity) { m.parent = p }

func (m *msg) addEnum(e Enum) {
	e.setParent(m)
	m.enums = append(m.enums, e)
}

func (m *msg) addMessage(sm Message) {
	sm.setParent(m)
	m.msgs = append(m.msgs, sm)
}

func (m *msg) addField(f Field) {
	f.setMessage(m)
	m.fields = append(m.fields, f)
}

func (m *msg) addOneOf(o OneOf) {
	o.setMessage(m)
	m.oneofs = append(m.oneofs, o)
}

func (m *msg) addMapEntry(me Message) {
	me.setParent(m)
	m.maps = append(m.maps, me)
}

func (m *msg) addDependent(message Message) {
	m.dependents = append(m.dependents, message)
}

func (m *msg) addDependency(message Message) {
	if message != nil {
		m.dependencies = append(m.dependencies, message)
	}
}

func (m *msg) childAtPath(path []int32) Entity {
	switch {
	case len(path) == 0:
		return m
	case len(path)%2 != 0:
		return nil
	}

	var child Entity
	switch path[0] {
	case messageTypeFieldPath:
		child = m.fields[path[1]]
	case messageTypeNestedTypePath:
		child = m.preservedMsgs[path[1]]
	case messageTypeEnumTypePath:
		child = m.enums[path[1]]
	case messageTypeOneofDeclPath:
		child = m.oneofs[path[1]]
	default:
		return nil
	}

	return child.childAtPath(path[2:])
}

func (m *msg) addSourceCodeInfo(info SourceCodeInfo) { m.info = info }

func messageSetToSlice(name string, set map[string]Message) []Message {
	dependents := make([]Message, 0, len(set))

	for fqn, d := range set {
		if fqn != name {
			dependents = append(dependents, d)
		}
	}

	return dependents
}

var _ Message = (*msg)(nil)
