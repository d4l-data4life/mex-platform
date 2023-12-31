// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.0
// source: shared/fields/fields.proto

package fields

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IndexDef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MultiValued bool         `protobuf:"varint,1,opt,name=multi_valued,json=multiValued,proto3" json:"multi_valued,omitempty"`
	Ext         []*anypb.Any `protobuf:"bytes,7,rep,name=ext,proto3" json:"ext,omitempty"`
}

func (x *IndexDef) Reset() {
	*x = IndexDef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_fields_fields_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexDef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexDef) ProtoMessage() {}

func (x *IndexDef) ProtoReflect() protoreflect.Message {
	mi := &file_shared_fields_fields_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexDef.ProtoReflect.Descriptor instead.
func (*IndexDef) Descriptor() ([]byte, []int) {
	return file_shared_fields_fields_proto_rawDescGZIP(), []int{0}
}

func (x *IndexDef) GetMultiValued() bool {
	if x != nil {
		return x.MultiValued
	}
	return false
}

func (x *IndexDef) GetExt() []*anypb.Any {
	if x != nil {
		return x.Ext
	}
	return nil
}

type FieldDef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Kind      string    `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	DisplayId string    `protobuf:"bytes,3,opt,name=display_id,json=displayId,proto3" json:"display_id,omitempty"`
	IndexDef  *IndexDef `protobuf:"bytes,4,opt,name=index_def,json=indexDef,proto3" json:"index_def,omitempty"`
}

func (x *FieldDef) Reset() {
	*x = FieldDef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_fields_fields_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldDef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldDef) ProtoMessage() {}

func (x *FieldDef) ProtoReflect() protoreflect.Message {
	mi := &file_shared_fields_fields_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldDef.ProtoReflect.Descriptor instead.
func (*FieldDef) Descriptor() ([]byte, []int) {
	return file_shared_fields_fields_proto_rawDescGZIP(), []int{1}
}

func (x *FieldDef) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FieldDef) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *FieldDef) GetDisplayId() string {
	if x != nil {
		return x.DisplayId
	}
	return ""
}

func (x *FieldDef) GetIndexDef() *IndexDef {
	if x != nil {
		return x.IndexDef
	}
	return nil
}

type FieldDefList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FieldDefs []*FieldDef `protobuf:"bytes,1,rep,name=field_defs,json=fieldDefs,proto3" json:"field_defs,omitempty"`
}

func (x *FieldDefList) Reset() {
	*x = FieldDefList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_fields_fields_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldDefList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldDefList) ProtoMessage() {}

func (x *FieldDefList) ProtoReflect() protoreflect.Message {
	mi := &file_shared_fields_fields_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldDefList.ProtoReflect.Descriptor instead.
func (*FieldDefList) Descriptor() ([]byte, []int) {
	return file_shared_fields_fields_proto_rawDescGZIP(), []int{2}
}

func (x *FieldDefList) GetFieldDefs() []*FieldDef {
	if x != nil {
		return x.FieldDefs
	}
	return nil
}

type IndexDefExtHierarchy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CodeSystemNameOrNodeEntityType string `protobuf:"bytes,1,opt,name=code_system_name_or_node_entity_type,json=codeSystemNameOrNodeEntityType,proto3" json:"code_system_name_or_node_entity_type,omitempty"`
	LinkFieldName                  string `protobuf:"bytes,4,opt,name=link_field_name,json=linkFieldName,proto3" json:"link_field_name,omitempty"`
	DisplayFieldName               string `protobuf:"bytes,5,opt,name=display_field_name,json=displayFieldName,proto3" json:"display_field_name,omitempty"`
}

func (x *IndexDefExtHierarchy) Reset() {
	*x = IndexDefExtHierarchy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_fields_fields_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexDefExtHierarchy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexDefExtHierarchy) ProtoMessage() {}

func (x *IndexDefExtHierarchy) ProtoReflect() protoreflect.Message {
	mi := &file_shared_fields_fields_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexDefExtHierarchy.ProtoReflect.Descriptor instead.
func (*IndexDefExtHierarchy) Descriptor() ([]byte, []int) {
	return file_shared_fields_fields_proto_rawDescGZIP(), []int{3}
}

func (x *IndexDefExtHierarchy) GetCodeSystemNameOrNodeEntityType() string {
	if x != nil {
		return x.CodeSystemNameOrNodeEntityType
	}
	return ""
}

func (x *IndexDefExtHierarchy) GetLinkFieldName() string {
	if x != nil {
		return x.LinkFieldName
	}
	return ""
}

func (x *IndexDefExtHierarchy) GetDisplayFieldName() string {
	if x != nil {
		return x.DisplayFieldName
	}
	return ""
}

type IndexDefExtLink struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RelationType       string   `protobuf:"bytes,1,opt,name=relation_type,json=relationType,proto3" json:"relation_type,omitempty"`
	LinkedTargetFields []string `protobuf:"bytes,4,rep,name=linked_target_fields,json=linkedTargetFields,proto3" json:"linked_target_fields,omitempty"`
}

func (x *IndexDefExtLink) Reset() {
	*x = IndexDefExtLink{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_fields_fields_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexDefExtLink) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexDefExtLink) ProtoMessage() {}

func (x *IndexDefExtLink) ProtoReflect() protoreflect.Message {
	mi := &file_shared_fields_fields_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexDefExtLink.ProtoReflect.Descriptor instead.
func (*IndexDefExtLink) Descriptor() ([]byte, []int) {
	return file_shared_fields_fields_proto_rawDescGZIP(), []int{4}
}

func (x *IndexDefExtLink) GetRelationType() string {
	if x != nil {
		return x.RelationType
	}
	return ""
}

func (x *IndexDefExtLink) GetLinkedTargetFields() []string {
	if x != nil {
		return x.LinkedTargetFields
	}
	return nil
}

type IndexDefExtCoding struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CodingsetNames []string `protobuf:"bytes,1,rep,name=codingset_names,json=codingsetNames,proto3" json:"codingset_names,omitempty"`
}

func (x *IndexDefExtCoding) Reset() {
	*x = IndexDefExtCoding{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_fields_fields_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexDefExtCoding) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexDefExtCoding) ProtoMessage() {}

func (x *IndexDefExtCoding) ProtoReflect() protoreflect.Message {
	mi := &file_shared_fields_fields_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexDefExtCoding.ProtoReflect.Descriptor instead.
func (*IndexDefExtCoding) Descriptor() ([]byte, []int) {
	return file_shared_fields_fields_proto_rawDescGZIP(), []int{5}
}

func (x *IndexDefExtCoding) GetCodingsetNames() []string {
	if x != nil {
		return x.CodingsetNames
	}
	return nil
}

var File_shared_fields_fields_proto protoreflect.FileDescriptor

var file_shared_fields_fields_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x2f,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6d, 0x65,
	0x78, 0x2e, 0x76, 0x30, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x55, 0x0a, 0x08, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65, 0x66, 0x12, 0x21, 0x0a, 0x0c, 0x6d,
	0x75, 0x6c, 0x74, 0x69, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x64, 0x12, 0x26,
	0x0a, 0x03, 0x65, 0x78, 0x74, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x03, 0x65, 0x78, 0x74, 0x22, 0xc7, 0x01, 0x0a, 0x08, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x44, 0x65, 0x66, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x64,
	0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x09, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x5f, 0x64, 0x65, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x6d, 0x65, 0x78, 0x2e, 0x76, 0x30, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65, 0x66, 0x52,
	0x08, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65, 0x66, 0x3a, 0x45, 0x92, 0x41, 0x42, 0x32, 0x40,
	0x7b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0x2c,
	0x22, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x3a, 0x22, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x2c,
	0x22, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65, 0x66, 0x22, 0x3a, 0x7b, 0x22, 0x6d, 0x75, 0x6c,
	0x74, 0x69, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x64, 0x22, 0x3a, 0x74, 0x72, 0x75, 0x65, 0x7d, 0x7d,
	0x22, 0x3f, 0x0a, 0x0c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x44, 0x65, 0x66, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x2f, 0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x64, 0x65, 0x66, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x76, 0x30, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x44, 0x65, 0x66, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x44, 0x65, 0x66,
	0x73, 0x22, 0xba, 0x01, 0x0a, 0x14, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65, 0x66, 0x45, 0x78,
	0x74, 0x48, 0x69, 0x65, 0x72, 0x61, 0x72, 0x63, 0x68, 0x79, 0x12, 0x4c, 0x0a, 0x24, 0x63, 0x6f,
	0x64, 0x65, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x6f,
	0x72, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1e, 0x63, 0x6f, 0x64, 0x65, 0x53, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x4f, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x6c, 0x69, 0x6e, 0x6b,
	0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x6c, 0x69, 0x6e, 0x6b, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x2c, 0x0a, 0x12, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x64, 0x69,
	0x73, 0x70, 0x6c, 0x61, 0x79, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x68,
	0x0a, 0x0f, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x65, 0x66, 0x45, 0x78, 0x74, 0x4c, 0x69, 0x6e,
	0x6b, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64,
	0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x12, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x3c, 0x0a, 0x11, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x44, 0x65, 0x66, 0x45, 0x78, 0x74, 0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x27, 0x0a,
	0x0f, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x65,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x42, 0x47, 0x5a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x73, 0x75, 0x6e, 0x64, 0x68, 0x65, 0x69, 0x74, 0x73,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72, 0x6b, 0x69, 0x2d, 0x6d, 0x65, 0x78, 0x2d, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x78, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x3b, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shared_fields_fields_proto_rawDescOnce sync.Once
	file_shared_fields_fields_proto_rawDescData = file_shared_fields_fields_proto_rawDesc
)

func file_shared_fields_fields_proto_rawDescGZIP() []byte {
	file_shared_fields_fields_proto_rawDescOnce.Do(func() {
		file_shared_fields_fields_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_fields_fields_proto_rawDescData)
	})
	return file_shared_fields_fields_proto_rawDescData
}

var file_shared_fields_fields_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_shared_fields_fields_proto_goTypes = []interface{}{
	(*IndexDef)(nil),             // 0: mex.v0.IndexDef
	(*FieldDef)(nil),             // 1: mex.v0.FieldDef
	(*FieldDefList)(nil),         // 2: mex.v0.FieldDefList
	(*IndexDefExtHierarchy)(nil), // 3: mex.v0.IndexDefExtHierarchy
	(*IndexDefExtLink)(nil),      // 4: mex.v0.IndexDefExtLink
	(*IndexDefExtCoding)(nil),    // 5: mex.v0.IndexDefExtCoding
	(*anypb.Any)(nil),            // 6: google.protobuf.Any
}
var file_shared_fields_fields_proto_depIdxs = []int32{
	6, // 0: mex.v0.IndexDef.ext:type_name -> google.protobuf.Any
	0, // 1: mex.v0.FieldDef.index_def:type_name -> mex.v0.IndexDef
	1, // 2: mex.v0.FieldDefList.field_defs:type_name -> mex.v0.FieldDef
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_shared_fields_fields_proto_init() }
func file_shared_fields_fields_proto_init() {
	if File_shared_fields_fields_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_fields_fields_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexDef); i {
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
		file_shared_fields_fields_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldDef); i {
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
		file_shared_fields_fields_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldDefList); i {
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
		file_shared_fields_fields_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexDefExtHierarchy); i {
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
		file_shared_fields_fields_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexDefExtLink); i {
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
		file_shared_fields_fields_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexDefExtCoding); i {
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
			RawDescriptor: file_shared_fields_fields_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_fields_fields_proto_goTypes,
		DependencyIndexes: file_shared_fields_fields_proto_depIdxs,
		MessageInfos:      file_shared_fields_fields_proto_msgTypes,
	}.Build()
	File_shared_fields_fields_proto = out.File
	file_shared_fields_fields_proto_rawDesc = nil
	file_shared_fields_fields_proto_goTypes = nil
	file_shared_fields_fields_proto_depIdxs = nil
}
