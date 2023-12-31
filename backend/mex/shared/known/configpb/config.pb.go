// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.0
// source: d4l/config.proto

package configpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Options struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If the value of the effective environment variable is empty or the variable is not set at all,
	// this value is used (as if the variable would have had this value).
	// Not specifying this option is equivalent to setting it to the empty string.
	// If a default is specified, the field cannot be made secret.
	Default string `protobuf:"bytes,1,opt,name=default,proto3" json:"default,omitempty"`
	// Indicates that the config field is confidential.
	// When printed, the config field is redacted.
	// When translated to K8s descriptors, the field becomes part
	// of a K8s Secret (and not a K8s ConfigMap).
	// A secret field cannot have a default value.
	Secret bool `protobuf:"varint,2,opt,name=secret,proto3" json:"secret,omitempty"`
	// Can be used to explicitly specifiy the environment variable name,
	// instead of the one derived from the field path in the Golang struct.
	Env string `protobuf:"bytes,3,opt,name=env,proto3" json:"env,omitempty"`
	// If true, the field is treated as if absent.
	// No translation or parsing takes place.
	Ignore bool `protobuf:"varint,4,opt,name=ignore,proto3" json:"ignore,omitempty"`
}

func (x *Options) Reset() {
	*x = Options{}
	if protoimpl.UnsafeEnabled {
		mi := &file_d4l_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Options) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Options) ProtoMessage() {}

func (x *Options) ProtoReflect() protoreflect.Message {
	mi := &file_d4l_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Options.ProtoReflect.Descriptor instead.
func (*Options) Descriptor() ([]byte, []int) {
	return file_d4l_config_proto_rawDescGZIP(), []int{0}
}

func (x *Options) GetDefault() string {
	if x != nil {
		return x.Default
	}
	return ""
}

func (x *Options) GetSecret() bool {
	if x != nil {
		return x.Secret
	}
	return false
}

func (x *Options) GetEnv() string {
	if x != nil {
		return x.Env
	}
	return ""
}

func (x *Options) GetIgnore() bool {
	if x != nil {
		return x.Ignore
	}
	return false
}

type Descriptor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title       string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Summary     string   `protobuf:"bytes,2,opt,name=summary,proto3" json:"summary,omitempty"`
	Description []string `protobuf:"bytes,3,rep,name=description,proto3" json:"description,omitempty"`
}

func (x *Descriptor) Reset() {
	*x = Descriptor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_d4l_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Descriptor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Descriptor) ProtoMessage() {}

func (x *Descriptor) ProtoReflect() protoreflect.Message {
	mi := &file_d4l_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Descriptor.ProtoReflect.Descriptor instead.
func (*Descriptor) Descriptor() ([]byte, []int) {
	return file_d4l_config_proto_rawDescGZIP(), []int{1}
}

func (x *Descriptor) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Descriptor) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *Descriptor) GetDescription() []string {
	if x != nil {
		return x.Description
	}
	return nil
}

type Export struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If true does not get turned into a field in a ConfigMap or Secret.
	Ignore bool `protobuf:"varint,1,opt,name=ignore,proto3" json:"ignore,omitempty"`
	// Can be used to set the source environment variable (otherwise the generated name is used).
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *Export) Reset() {
	*x = Export{}
	if protoimpl.UnsafeEnabled {
		mi := &file_d4l_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Export) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Export) ProtoMessage() {}

func (x *Export) ProtoReflect() protoreflect.Message {
	mi := &file_d4l_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Export.ProtoReflect.Descriptor instead.
func (*Export) Descriptor() ([]byte, []int) {
	return file_d4l_config_proto_rawDescGZIP(), []int{2}
}

func (x *Export) GetIgnore() bool {
	if x != nil {
		return x.Ignore
	}
	return false
}

func (x *Export) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

var file_d4l_config_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*Options)(nil),
		Field:         20000,
		Name:          "d4l.cfg.opts",
		Tag:           "bytes,20000,opt,name=opts",
		Filename:      "d4l/config.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*Descriptor)(nil),
		Field:         20001,
		Name:          "d4l.cfg.desc",
		Tag:           "bytes,20001,opt,name=desc",
		Filename:      "d4l/config.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         20002,
		Name:          "d4l.cfg.tags",
		Tag:           "bytes,20002,rep,name=tags",
		Filename:      "d4l/config.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*Export)(nil),
		Field:         20003,
		Name:          "d4l.cfg.k8s",
		Tag:           "bytes,20003,opt,name=k8s",
		Filename:      "d4l/config.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*Descriptor)(nil),
		Field:         20002,
		Name:          "d4l.cfg.mdesc",
		Tag:           "bytes,20002,opt,name=mdesc",
		Filename:      "d4l/config.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         20003,
		Name:          "d4l.cfg.mtags",
		Tag:           "bytes,20003,rep,name=mtags",
		Filename:      "d4l/config.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50000,
		Name:          "d4l.cfg.main_message",
		Tag:           "bytes,50000,opt,name=main_message",
		Filename:      "d4l/config.proto",
	},
}

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional d4l.cfg.Options opts = 20000;
	E_Opts = &file_d4l_config_proto_extTypes[0]
	// optional d4l.cfg.Descriptor desc = 20001;
	E_Desc = &file_d4l_config_proto_extTypes[1]
	// repeated string tags = 20002;
	E_Tags = &file_d4l_config_proto_extTypes[2]
	// optional d4l.cfg.Export k8s = 20003;
	E_K8S = &file_d4l_config_proto_extTypes[3]
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional d4l.cfg.Descriptor mdesc = 20002;
	E_Mdesc = &file_d4l_config_proto_extTypes[4]
	// repeated string mtags = 20003;
	E_Mtags = &file_d4l_config_proto_extTypes[5]
)

// Extension fields to descriptorpb.FileOptions.
var (
	// optional string main_message = 50000;
	E_MainMessage = &file_d4l_config_proto_extTypes[6]
)

var File_d4l_config_proto protoreflect.FileDescriptor

var file_d4l_config_proto_rawDesc = []byte{
	0x0a, 0x10, 0x64, 0x34, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x64, 0x34, 0x6c, 0x2e, 0x63, 0x66, 0x67, 0x1a, 0x20, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a,
	0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x66, 0x61,
	0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e,
	0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e, 0x76, 0x12, 0x16, 0x0a, 0x06,
	0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x67,
	0x6e, 0x6f, 0x72, 0x65, 0x22, 0x5e, 0x0a, 0x0a, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x38, 0x0a, 0x06, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x3a, 0x45,
	0x0a, 0x04, 0x6f, 0x70, 0x74, 0x73, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa0, 0x9c, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x64, 0x34, 0x6c, 0x2e, 0x63, 0x66, 0x67, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x04, 0x6f, 0x70, 0x74, 0x73, 0x3a, 0x48, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x1d, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa1, 0x9c, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x63, 0x66, 0x67, 0x2e, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x3a,
	0x33, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0x9c, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x3a, 0x42, 0x0a, 0x03, 0x6b, 0x38, 0x73, 0x12, 0x1d, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa3, 0x9c, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x63, 0x66, 0x67, 0x2e, 0x45, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x52, 0x03, 0x6b, 0x38, 0x73, 0x3a, 0x4c, 0x0a, 0x05, 0x6d, 0x64, 0x65, 0x73,
	0x63, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xa2, 0x9c, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x64, 0x34, 0x6c,
	0x2e, 0x63, 0x66, 0x67, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52,
	0x05, 0x6d, 0x64, 0x65, 0x73, 0x63, 0x3a, 0x37, 0x0a, 0x05, 0x6d, 0x74, 0x61, 0x67, 0x73, 0x12,
	0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xa3, 0x9c, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x74, 0x61, 0x67, 0x73, 0x3a,
	0x41, 0x0a, 0x0c, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x61, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x67, 0x65, 0x73, 0x75, 0x6e, 0x64, 0x68, 0x65, 0x69, 0x74, 0x73, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x2f, 0x72, 0x6b, 0x69, 0x2d, 0x6d, 0x65, 0x78, 0x2d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2f, 0x6d, 0x65, 0x78, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x6b, 0x6e,
	0x6f, 0x77, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_d4l_config_proto_rawDescOnce sync.Once
	file_d4l_config_proto_rawDescData = file_d4l_config_proto_rawDesc
)

func file_d4l_config_proto_rawDescGZIP() []byte {
	file_d4l_config_proto_rawDescOnce.Do(func() {
		file_d4l_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_d4l_config_proto_rawDescData)
	})
	return file_d4l_config_proto_rawDescData
}

var file_d4l_config_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_d4l_config_proto_goTypes = []interface{}{
	(*Options)(nil),                     // 0: d4l.cfg.Options
	(*Descriptor)(nil),                  // 1: d4l.cfg.Descriptor
	(*Export)(nil),                      // 2: d4l.cfg.Export
	(*descriptorpb.FieldOptions)(nil),   // 3: google.protobuf.FieldOptions
	(*descriptorpb.MessageOptions)(nil), // 4: google.protobuf.MessageOptions
	(*descriptorpb.FileOptions)(nil),    // 5: google.protobuf.FileOptions
}
var file_d4l_config_proto_depIdxs = []int32{
	3,  // 0: d4l.cfg.opts:extendee -> google.protobuf.FieldOptions
	3,  // 1: d4l.cfg.desc:extendee -> google.protobuf.FieldOptions
	3,  // 2: d4l.cfg.tags:extendee -> google.protobuf.FieldOptions
	3,  // 3: d4l.cfg.k8s:extendee -> google.protobuf.FieldOptions
	4,  // 4: d4l.cfg.mdesc:extendee -> google.protobuf.MessageOptions
	4,  // 5: d4l.cfg.mtags:extendee -> google.protobuf.MessageOptions
	5,  // 6: d4l.cfg.main_message:extendee -> google.protobuf.FileOptions
	0,  // 7: d4l.cfg.opts:type_name -> d4l.cfg.Options
	1,  // 8: d4l.cfg.desc:type_name -> d4l.cfg.Descriptor
	2,  // 9: d4l.cfg.k8s:type_name -> d4l.cfg.Export
	1,  // 10: d4l.cfg.mdesc:type_name -> d4l.cfg.Descriptor
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	7,  // [7:11] is the sub-list for extension type_name
	0,  // [0:7] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_d4l_config_proto_init() }
func file_d4l_config_proto_init() {
	if File_d4l_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_d4l_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Options); i {
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
		file_d4l_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Descriptor); i {
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
		file_d4l_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Export); i {
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
			RawDescriptor: file_d4l_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 7,
			NumServices:   0,
		},
		GoTypes:           file_d4l_config_proto_goTypes,
		DependencyIndexes: file_d4l_config_proto_depIdxs,
		MessageInfos:      file_d4l_config_proto_msgTypes,
		ExtensionInfos:    file_d4l_config_proto_extTypes,
	}.Build()
	File_d4l_config_proto = out.File
	file_d4l_config_proto_rawDesc = nil
	file_d4l_config_proto_goTypes = nil
	file_d4l_config_proto_depIdxs = nil
}
