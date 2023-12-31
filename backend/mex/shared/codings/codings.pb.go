// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.0
// source: shared/codings/codings.proto

package codings

import (
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

type BlobStoreCodingsetSourceConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlobName string `protobuf:"bytes,1,opt,name=blob_name,json=blobName,proto3" json:"blob_name,omitempty"`
	BlobType string `protobuf:"bytes,2,opt,name=blob_type,json=blobType,proto3" json:"blob_type,omitempty"`
}

func (x *BlobStoreCodingsetSourceConfig) Reset() {
	*x = BlobStoreCodingsetSourceConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_codings_codings_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlobStoreCodingsetSourceConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlobStoreCodingsetSourceConfig) ProtoMessage() {}

func (x *BlobStoreCodingsetSourceConfig) ProtoReflect() protoreflect.Message {
	mi := &file_shared_codings_codings_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlobStoreCodingsetSourceConfig.ProtoReflect.Descriptor instead.
func (*BlobStoreCodingsetSourceConfig) Descriptor() ([]byte, []int) {
	return file_shared_codings_codings_proto_rawDescGZIP(), []int{0}
}

func (x *BlobStoreCodingsetSourceConfig) GetBlobName() string {
	if x != nil {
		return x.BlobName
	}
	return ""
}

func (x *BlobStoreCodingsetSourceConfig) GetBlobType() string {
	if x != nil {
		return x.BlobType
	}
	return ""
}

type CodingsetSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Config *anypb.Any `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *CodingsetSource) Reset() {
	*x = CodingsetSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_codings_codings_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CodingsetSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CodingsetSource) ProtoMessage() {}

func (x *CodingsetSource) ProtoReflect() protoreflect.Message {
	mi := &file_shared_codings_codings_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CodingsetSource.ProtoReflect.Descriptor instead.
func (*CodingsetSource) Descriptor() ([]byte, []int) {
	return file_shared_codings_codings_proto_rawDescGZIP(), []int{1}
}

func (x *CodingsetSource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CodingsetSource) GetConfig() *anypb.Any {
	if x != nil {
		return x.Config
	}
	return nil
}

type CodingsetSources struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CodingsetSources []*CodingsetSource `protobuf:"bytes,1,rep,name=codingset_sources,json=codingsetSources,proto3" json:"codingset_sources,omitempty"`
}

func (x *CodingsetSources) Reset() {
	*x = CodingsetSources{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_codings_codings_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CodingsetSources) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CodingsetSources) ProtoMessage() {}

func (x *CodingsetSources) ProtoReflect() protoreflect.Message {
	mi := &file_shared_codings_codings_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CodingsetSources.ProtoReflect.Descriptor instead.
func (*CodingsetSources) Descriptor() ([]byte, []int) {
	return file_shared_codings_codings_proto_rawDescGZIP(), []int{2}
}

func (x *CodingsetSources) GetCodingsetSources() []*CodingsetSource {
	if x != nil {
		return x.CodingsetSources
	}
	return nil
}

var File_shared_codings_codings_proto protoreflect.FileDescriptor

var file_shared_codings_codings_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x2f, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x1a,
	0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5a, 0x0a, 0x1e, 0x42, 0x6c,
	0x6f, 0x62, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x65, 0x74,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1b, 0x0a, 0x09,
	0x62, 0x6c, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x62, 0x6c, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6c, 0x6f,
	0x62, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x6c,
	0x6f, 0x62, 0x54, 0x79, 0x70, 0x65, 0x22, 0x53, 0x0a, 0x0f, 0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67,
	0x73, 0x65, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a,
	0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x41, 0x6e, 0x79, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x61, 0x0a, 0x10, 0x43,
	0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x65, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12,
	0x4d, 0x0a, 0x11, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x65, 0x74, 0x5f, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x64, 0x34, 0x6c,
	0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x43, 0x6f, 0x64,
	0x69, 0x6e, 0x67, 0x73, 0x65, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x10, 0x63, 0x6f,
	0x64, 0x69, 0x6e, 0x67, 0x73, 0x65, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x42, 0x49,
	0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x73,
	0x75, 0x6e, 0x64, 0x68, 0x65, 0x69, 0x74, 0x73, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72, 0x6b,
	0x69, 0x2d, 0x6d, 0x65, 0x78, 0x2d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d,
	0x65, 0x78, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67,
	0x73, 0x3b, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_shared_codings_codings_proto_rawDescOnce sync.Once
	file_shared_codings_codings_proto_rawDescData = file_shared_codings_codings_proto_rawDesc
)

func file_shared_codings_codings_proto_rawDescGZIP() []byte {
	file_shared_codings_codings_proto_rawDescOnce.Do(func() {
		file_shared_codings_codings_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_codings_codings_proto_rawDescData)
	})
	return file_shared_codings_codings_proto_rawDescData
}

var file_shared_codings_codings_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_shared_codings_codings_proto_goTypes = []interface{}{
	(*BlobStoreCodingsetSourceConfig)(nil), // 0: d4l.mex.codings.BlobStoreCodingsetSourceConfig
	(*CodingsetSource)(nil),                // 1: d4l.mex.codings.CodingsetSource
	(*CodingsetSources)(nil),               // 2: d4l.mex.codings.CodingsetSources
	(*anypb.Any)(nil),                      // 3: google.protobuf.Any
}
var file_shared_codings_codings_proto_depIdxs = []int32{
	3, // 0: d4l.mex.codings.CodingsetSource.config:type_name -> google.protobuf.Any
	1, // 1: d4l.mex.codings.CodingsetSources.codingset_sources:type_name -> d4l.mex.codings.CodingsetSource
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_shared_codings_codings_proto_init() }
func file_shared_codings_codings_proto_init() {
	if File_shared_codings_codings_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_codings_codings_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlobStoreCodingsetSourceConfig); i {
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
		file_shared_codings_codings_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CodingsetSource); i {
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
		file_shared_codings_codings_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CodingsetSources); i {
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
			RawDescriptor: file_shared_codings_codings_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_codings_codings_proto_goTypes,
		DependencyIndexes: file_shared_codings_codings_proto_depIdxs,
		MessageInfos:      file_shared_codings_codings_proto_msgTypes,
	}.Build()
	File_shared_codings_codings_proto = out.File
	file_shared_codings_codings_proto_rawDesc = nil
	file_shared_codings_codings_proto_goTypes = nil
	file_shared_codings_codings_proto_depIdxs = nil
}
