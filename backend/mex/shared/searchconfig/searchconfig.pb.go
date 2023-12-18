// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.0
// source: shared/searchconfig/searchconfig.proto

package searchconfig

import (
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

type SearchConfigObject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type   string   `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *SearchConfigObject) Reset() {
	*x = SearchConfigObject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_searchconfig_searchconfig_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchConfigObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchConfigObject) ProtoMessage() {}

func (x *SearchConfigObject) ProtoReflect() protoreflect.Message {
	mi := &file_shared_searchconfig_searchconfig_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchConfigObject.ProtoReflect.Descriptor instead.
func (*SearchConfigObject) Descriptor() ([]byte, []int) {
	return file_shared_searchconfig_searchconfig_proto_rawDescGZIP(), []int{0}
}

func (x *SearchConfigObject) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SearchConfigObject) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *SearchConfigObject) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type SearchConfigList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SearchConfigs []*SearchConfigObject `protobuf:"bytes,1,rep,name=search_configs,json=searchConfigs,proto3" json:"search_configs,omitempty"`
}

func (x *SearchConfigList) Reset() {
	*x = SearchConfigList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_searchconfig_searchconfig_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchConfigList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchConfigList) ProtoMessage() {}

func (x *SearchConfigList) ProtoReflect() protoreflect.Message {
	mi := &file_shared_searchconfig_searchconfig_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchConfigList.ProtoReflect.Descriptor instead.
func (*SearchConfigList) Descriptor() ([]byte, []int) {
	return file_shared_searchconfig_searchconfig_proto_rawDescGZIP(), []int{1}
}

func (x *SearchConfigList) GetSearchConfigs() []*SearchConfigObject {
	if x != nil {
		return x.SearchConfigs
	}
	return nil
}

var File_shared_searchconfig_searchconfig_proto protoreflect.FileDescriptor

var file_shared_searchconfig_searchconfig_proto_rawDesc = []byte{
	0x0a, 0x26, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65,
	0x78, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x54,
	0x0a, 0x12, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x22, 0x63, 0x0a, 0x10, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x4f, 0x0a, 0x0e, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x28, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x0d, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x42, 0x53, 0x5a, 0x51, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x73, 0x75, 0x6e, 0x64, 0x68, 0x65,
	0x69, 0x74, 0x73, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72, 0x6b, 0x69, 0x2d, 0x6d, 0x65, 0x78,
	0x2d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x78, 0x2f, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x3b, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shared_searchconfig_searchconfig_proto_rawDescOnce sync.Once
	file_shared_searchconfig_searchconfig_proto_rawDescData = file_shared_searchconfig_searchconfig_proto_rawDesc
)

func file_shared_searchconfig_searchconfig_proto_rawDescGZIP() []byte {
	file_shared_searchconfig_searchconfig_proto_rawDescOnce.Do(func() {
		file_shared_searchconfig_searchconfig_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_searchconfig_searchconfig_proto_rawDescData)
	})
	return file_shared_searchconfig_searchconfig_proto_rawDescData
}

var file_shared_searchconfig_searchconfig_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_shared_searchconfig_searchconfig_proto_goTypes = []interface{}{
	(*SearchConfigObject)(nil), // 0: d4l.mex.searchconfig.SearchConfigObject
	(*SearchConfigList)(nil),   // 1: d4l.mex.searchconfig.SearchConfigList
}
var file_shared_searchconfig_searchconfig_proto_depIdxs = []int32{
	0, // 0: d4l.mex.searchconfig.SearchConfigList.search_configs:type_name -> d4l.mex.searchconfig.SearchConfigObject
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_shared_searchconfig_searchconfig_proto_init() }
func file_shared_searchconfig_searchconfig_proto_init() {
	if File_shared_searchconfig_searchconfig_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_searchconfig_searchconfig_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchConfigObject); i {
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
		file_shared_searchconfig_searchconfig_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchConfigList); i {
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
			RawDescriptor: file_shared_searchconfig_searchconfig_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_searchconfig_searchconfig_proto_goTypes,
		DependencyIndexes: file_shared_searchconfig_searchconfig_proto_depIdxs,
		MessageInfos:      file_shared_searchconfig_searchconfig_proto_msgTypes,
	}.Build()
	File_shared_searchconfig_searchconfig_proto = out.File
	file_shared_searchconfig_searchconfig_proto_rawDesc = nil
	file_shared_searchconfig_searchconfig_proto_goTypes = nil
	file_shared_searchconfig_searchconfig_proto_depIdxs = nil
}