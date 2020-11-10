// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: dynamo.proto

package dynamo

import (
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type DynamoMessageOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Disables generation for this message
	Disabled bool `protobuf:"varint,1,opt,name=disabled,proto3" json:"disabled,omitempty"`
	// A compound value that can be used as a Partition Key.
	Partition *Key `protobuf:"bytes,2,opt,name=partition,proto3" json:"partition,omitempty"`
	// A compound value that can be used as a Sort Key.
	Sort *Key `protobuf:"bytes,3,opt,name=sort,proto3" json:"sort,omitempty"`
	// A list of compound values that can be set from other fields.
	CompoundField []*Key `protobuf:"bytes,4,rep,name=compound_field,json=compoundField,proto3" json:"compound_field,omitempty"`
}

func (x *DynamoMessageOptions) Reset() {
	*x = DynamoMessageOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dynamo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DynamoMessageOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DynamoMessageOptions) ProtoMessage() {}

func (x *DynamoMessageOptions) ProtoReflect() protoreflect.Message {
	mi := &file_dynamo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DynamoMessageOptions.ProtoReflect.Descriptor instead.
func (*DynamoMessageOptions) Descriptor() ([]byte, []int) {
	return file_dynamo_proto_rawDescGZIP(), []int{0}
}

func (x *DynamoMessageOptions) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *DynamoMessageOptions) GetPartition() *Key {
	if x != nil {
		return x.Partition
	}
	return nil
}

func (x *DynamoMessageOptions) GetSort() *Key {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *DynamoMessageOptions) GetCompoundField() []*Key {
	if x != nil {
		return x.CompoundField
	}
	return nil
}

type Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	Prefix string   `protobuf:"bytes,3,opt,name=prefix,proto3" json:"prefix,omitempty"`
	// defaults to `:` if unset, only used if more than one field is set
	Separator string `protobuf:"bytes,4,opt,name=separator,proto3" json:"separator,omitempty"`
}

func (x *Key) Reset() {
	*x = Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dynamo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Key) ProtoMessage() {}

func (x *Key) ProtoReflect() protoreflect.Message {
	mi := &file_dynamo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Key.ProtoReflect.Descriptor instead.
func (*Key) Descriptor() ([]byte, []int) {
	return file_dynamo_proto_rawDescGZIP(), []int{1}
}

func (x *Key) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Key) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *Key) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *Key) GetSeparator() string {
	if x != nil {
		return x.Separator
	}
	return ""
}

type DynamoFieldOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Skip bool   `protobuf:"varint,1,opt,name=skip,proto3" json:"skip,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Type *Types `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *DynamoFieldOptions) Reset() {
	*x = DynamoFieldOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dynamo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DynamoFieldOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DynamoFieldOptions) ProtoMessage() {}

func (x *DynamoFieldOptions) ProtoReflect() protoreflect.Message {
	mi := &file_dynamo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DynamoFieldOptions.ProtoReflect.Descriptor instead.
func (*DynamoFieldOptions) Descriptor() ([]byte, []int) {
	return file_dynamo_proto_rawDescGZIP(), []int{2}
}

func (x *DynamoFieldOptions) GetSkip() bool {
	if x != nil {
		return x.Skip
	}
	return false
}

func (x *DynamoFieldOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DynamoFieldOptions) GetType() *Types {
	if x != nil {
		return x.Type
	}
	return nil
}

type Types struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Binary     bool `protobuf:"varint,100,opt,name=binary,proto3" json:"binary,omitempty"`
	Set        bool `protobuf:"varint,200,opt,name=set,proto3" json:"set,omitempty"`
	UnixSecond bool `protobuf:"varint,300,opt,name=unix_second,json=unixSecond,proto3" json:"unix_second,omitempty"`
	UnixMilli  bool `protobuf:"varint,301,opt,name=unix_milli,json=unixMilli,proto3" json:"unix_milli,omitempty"`
	UnixNano   bool `protobuf:"varint,302,opt,name=unix_nano,json=unixNano,proto3" json:"unix_nano,omitempty"`
}

func (x *Types) Reset() {
	*x = Types{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dynamo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Types) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Types) ProtoMessage() {}

func (x *Types) ProtoReflect() protoreflect.Message {
	mi := &file_dynamo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Types.ProtoReflect.Descriptor instead.
func (*Types) Descriptor() ([]byte, []int) {
	return file_dynamo_proto_rawDescGZIP(), []int{3}
}

func (x *Types) GetBinary() bool {
	if x != nil {
		return x.Binary
	}
	return false
}

func (x *Types) GetSet() bool {
	if x != nil {
		return x.Set
	}
	return false
}

func (x *Types) GetUnixSecond() bool {
	if x != nil {
		return x.UnixSecond
	}
	return false
}

func (x *Types) GetUnixMilli() bool {
	if x != nil {
		return x.UnixMilli
	}
	return false
}

func (x *Types) GetUnixNano() bool {
	if x != nil {
		return x.UnixNano
	}
	return false
}

var file_dynamo_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*DynamoMessageOptions)(nil),
		Field:         6000,
		Name:          "dynamo.msg",
		Tag:           "bytes,6000,opt,name=msg",
		Filename:      "dynamo.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*DynamoFieldOptions)(nil),
		Field:         6000,
		Name:          "dynamo.field",
		Tag:           "bytes,6000,opt,name=field",
		Filename:      "dynamo.proto",
	},
}

// Extension fields to descriptor.MessageOptions.
var (
	// optional dynamo.DynamoMessageOptions msg = 6000;
	E_Msg = &file_dynamo_proto_extTypes[0]
)

// Extension fields to descriptor.FieldOptions.
var (
	// optional dynamo.DynamoFieldOptions field = 6000;
	E_Field = &file_dynamo_proto_extTypes[1]
)

var File_dynamo_proto protoreflect.FileDescriptor

var file_dynamo_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb2, 0x01, 0x0a, 0x14, 0x44, 0x79, 0x6e,
	0x61, 0x6d, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x29, 0x0a,
	0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x09, 0x70,
	0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e,
	0x4b, 0x65, 0x79, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x32, 0x0a, 0x0e, 0x63, 0x6f, 0x6d,
	0x70, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0b, 0x2e, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x0d,
	0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x75, 0x6e, 0x64, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x22, 0x67, 0x0a,
	0x03, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x70, 0x61,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x70,
	0x61, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x5f, 0x0a, 0x12, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x6f,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x6b, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x73, 0x6b, 0x69, 0x70,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x92, 0x01, 0x0a, 0x05, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x18, 0x64, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x12, 0x11, 0x0a, 0x03, 0x73, 0x65, 0x74,
	0x18, 0xc8, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x73, 0x65, 0x74, 0x12, 0x20, 0x0a, 0x0b,
	0x75, 0x6e, 0x69, 0x78, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x18, 0xac, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0a, 0x75, 0x6e, 0x69, 0x78, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x18, 0xad, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x75, 0x6e, 0x69, 0x78, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x12, 0x1c,
	0x0a, 0x09, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0xae, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x3a, 0x50, 0x0a, 0x03,
	0x6d, 0x73, 0x67, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf0, 0x2e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x64, 0x79,
	0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x3a, 0x50,
	0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf0, 0x2e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2e, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70,
	0x71, 0x75, 0x65, 0x72, 0x6e, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65,
	0x6e, 0x2d, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x2f, 0x64, 0x79, 0x6e, 0x61, 0x6d, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dynamo_proto_rawDescOnce sync.Once
	file_dynamo_proto_rawDescData = file_dynamo_proto_rawDesc
)

func file_dynamo_proto_rawDescGZIP() []byte {
	file_dynamo_proto_rawDescOnce.Do(func() {
		file_dynamo_proto_rawDescData = protoimpl.X.CompressGZIP(file_dynamo_proto_rawDescData)
	})
	return file_dynamo_proto_rawDescData
}

var file_dynamo_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_dynamo_proto_goTypes = []interface{}{
	(*DynamoMessageOptions)(nil),      // 0: dynamo.DynamoMessageOptions
	(*Key)(nil),                       // 1: dynamo.Key
	(*DynamoFieldOptions)(nil),        // 2: dynamo.DynamoFieldOptions
	(*Types)(nil),                     // 3: dynamo.Types
	(*descriptor.MessageOptions)(nil), // 4: google.protobuf.MessageOptions
	(*descriptor.FieldOptions)(nil),   // 5: google.protobuf.FieldOptions
}
var file_dynamo_proto_depIdxs = []int32{
	1, // 0: dynamo.DynamoMessageOptions.partition:type_name -> dynamo.Key
	1, // 1: dynamo.DynamoMessageOptions.sort:type_name -> dynamo.Key
	1, // 2: dynamo.DynamoMessageOptions.compound_field:type_name -> dynamo.Key
	3, // 3: dynamo.DynamoFieldOptions.type:type_name -> dynamo.Types
	4, // 4: dynamo.msg:extendee -> google.protobuf.MessageOptions
	5, // 5: dynamo.field:extendee -> google.protobuf.FieldOptions
	0, // 6: dynamo.msg:type_name -> dynamo.DynamoMessageOptions
	2, // 7: dynamo.field:type_name -> dynamo.DynamoFieldOptions
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	6, // [6:8] is the sub-list for extension type_name
	4, // [4:6] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_dynamo_proto_init() }
func file_dynamo_proto_init() {
	if File_dynamo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dynamo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DynamoMessageOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_dynamo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Key); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_dynamo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DynamoFieldOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_dynamo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Types); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dynamo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_dynamo_proto_goTypes,
		DependencyIndexes: file_dynamo_proto_depIdxs,
		MessageInfos:      file_dynamo_proto_msgTypes,
		ExtensionInfos:    file_dynamo_proto_extTypes,
	}.Build()
	File_dynamo_proto = out.File
	file_dynamo_proto_rawDesc = nil
	file_dynamo_proto_goTypes = nil
	file_dynamo_proto_depIdxs = nil
}
