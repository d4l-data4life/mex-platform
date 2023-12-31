// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.0
// source: shared/solr/solr.proto

package solr

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

type Sorting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Axis  string `protobuf:"bytes,1,opt,name=axis,proto3" json:"axis,omitempty"`
	Order string `protobuf:"bytes,2,opt,name=order,proto3" json:"order,omitempty"`
}

func (x *Sorting) Reset() {
	*x = Sorting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sorting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sorting) ProtoMessage() {}

func (x *Sorting) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sorting.ProtoReflect.Descriptor instead.
func (*Sorting) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{0}
}

func (x *Sorting) GetAxis() string {
	if x != nil {
		return x.Axis
	}
	return ""
}

func (x *Sorting) GetOrder() string {
	if x != nil {
		return x.Order
	}
	return ""
}

type StringRange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Min string `protobuf:"bytes,1,opt,name=min,proto3" json:"min,omitempty"`
	Max string `protobuf:"bytes,2,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *StringRange) Reset() {
	*x = StringRange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringRange) ProtoMessage() {}

func (x *StringRange) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringRange.ProtoReflect.Descriptor instead.
func (*StringRange) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{1}
}

func (x *StringRange) GetMin() string {
	if x != nil {
		return x.Min
	}
	return ""
}

func (x *StringRange) GetMax() string {
	if x != nil {
		return x.Max
	}
	return ""
}

type AxisConstraint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type             string         `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Axis             string         `protobuf:"bytes,2,opt,name=axis,proto3" json:"axis,omitempty"`
	Values           []string       `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
	SingleNodeValues []string       `protobuf:"bytes,4,rep,name=single_node_values,json=singleNodeValues,proto3" json:"single_node_values,omitempty"`
	StringRanges     []*StringRange `protobuf:"bytes,5,rep,name=string_ranges,json=stringRanges,proto3" json:"string_ranges,omitempty"`
	CombineOperator  string         `protobuf:"bytes,6,opt,name=combine_operator,json=combineOperator,proto3" json:"combine_operator,omitempty"`
}

func (x *AxisConstraint) Reset() {
	*x = AxisConstraint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AxisConstraint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AxisConstraint) ProtoMessage() {}

func (x *AxisConstraint) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AxisConstraint.ProtoReflect.Descriptor instead.
func (*AxisConstraint) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{2}
}

func (x *AxisConstraint) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *AxisConstraint) GetAxis() string {
	if x != nil {
		return x.Axis
	}
	return ""
}

func (x *AxisConstraint) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *AxisConstraint) GetSingleNodeValues() []string {
	if x != nil {
		return x.SingleNodeValues
	}
	return nil
}

func (x *AxisConstraint) GetStringRanges() []*StringRange {
	if x != nil {
		return x.StringRanges
	}
	return nil
}

func (x *AxisConstraint) GetCombineOperator() string {
	if x != nil {
		return x.CombineOperator
	}
	return ""
}

type Facet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Axis     string `protobuf:"bytes,2,opt,name=axis,proto3" json:"axis,omitempty"`
	Limit    uint32 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset   uint32 `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
	StatName string `protobuf:"bytes,5,opt,name=stat_name,json=statName,proto3" json:"stat_name,omitempty"`
	StatOp   string `protobuf:"bytes,6,opt,name=stat_op,json=statOp,proto3" json:"stat_op,omitempty"`
}

func (x *Facet) Reset() {
	*x = Facet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Facet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Facet) ProtoMessage() {}

func (x *Facet) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Facet.ProtoReflect.Descriptor instead.
func (*Facet) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{3}
}

func (x *Facet) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Facet) GetAxis() string {
	if x != nil {
		return x.Axis
	}
	return ""
}

func (x *Facet) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *Facet) GetOffset() uint32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *Facet) GetStatName() string {
	if x != nil {
		return x.StatName
	}
	return ""
}

func (x *Facet) GetStatOp() string {
	if x != nil {
		return x.StatOp
	}
	return ""
}

type DocValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FieldName  string `protobuf:"bytes,1,opt,name=field_name,json=fieldName,proto3" json:"field_name,omitempty"`
	FieldValue string `protobuf:"bytes,2,opt,name=field_value,json=fieldValue,proto3" json:"field_value,omitempty"`
	Language   string `protobuf:"bytes,3,opt,name=language,proto3" json:"language,omitempty"`
}

func (x *DocValue) Reset() {
	*x = DocValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocValue) ProtoMessage() {}

func (x *DocValue) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocValue.ProtoReflect.Descriptor instead.
func (*DocValue) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{4}
}

func (x *DocValue) GetFieldName() string {
	if x != nil {
		return x.FieldName
	}
	return ""
}

func (x *DocValue) GetFieldValue() string {
	if x != nil {
		return x.FieldValue
	}
	return ""
}

func (x *DocValue) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

type DocItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemId     string      `protobuf:"bytes,1,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	EntityType string      `protobuf:"bytes,2,opt,name=entity_type,json=entityType,proto3" json:"entity_type,omitempty"`
	Values     []*DocValue `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *DocItem) Reset() {
	*x = DocItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocItem) ProtoMessage() {}

func (x *DocItem) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocItem.ProtoReflect.Descriptor instead.
func (*DocItem) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{5}
}

func (x *DocItem) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *DocItem) GetEntityType() string {
	if x != nil {
		return x.EntityType
	}
	return ""
}

func (x *DocItem) GetValues() []*DocValue {
	if x != nil {
		return x.Values
	}
	return nil
}

type HierarchyInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParentValue string `protobuf:"bytes,1,opt,name=parent_value,json=parentValue,proto3" json:"parent_value,omitempty"`
	Display     string `protobuf:"bytes,2,opt,name=display,proto3" json:"display,omitempty"`
	Depth       uint32 `protobuf:"varint,3,opt,name=depth,proto3" json:"depth,omitempty"`
}

func (x *HierarchyInfo) Reset() {
	*x = HierarchyInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HierarchyInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HierarchyInfo) ProtoMessage() {}

func (x *HierarchyInfo) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HierarchyInfo.ProtoReflect.Descriptor instead.
func (*HierarchyInfo) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{6}
}

func (x *HierarchyInfo) GetParentValue() string {
	if x != nil {
		return x.ParentValue
	}
	return ""
}

func (x *HierarchyInfo) GetDisplay() string {
	if x != nil {
		return x.Display
	}
	return ""
}

func (x *HierarchyInfo) GetDepth() uint32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

type FacetBucket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value         string     `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Count         uint32     `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	HierarchyInfo *anypb.Any `protobuf:"bytes,3,opt,name=hierarchyInfo,proto3" json:"hierarchyInfo,omitempty"`
}

func (x *FacetBucket) Reset() {
	*x = FacetBucket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FacetBucket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FacetBucket) ProtoMessage() {}

func (x *FacetBucket) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FacetBucket.ProtoReflect.Descriptor instead.
func (*FacetBucket) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{7}
}

func (x *FacetBucket) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *FacetBucket) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *FacetBucket) GetHierarchyInfo() *anypb.Any {
	if x != nil {
		return x.HierarchyInfo
	}
	return nil
}

type FacetResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type             string         `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Axis             string         `protobuf:"bytes,2,opt,name=axis,proto3" json:"axis,omitempty"`
	BucketNo         uint32         `protobuf:"varint,3,opt,name=bucketNo,proto3" json:"bucketNo,omitempty"`
	Buckets          []*FacetBucket `protobuf:"bytes,4,rep,name=buckets,proto3" json:"buckets,omitempty"`
	StatName         string         `protobuf:"bytes,5,opt,name=statName,proto3" json:"statName,omitempty"`
	StringStatResult string         `protobuf:"bytes,6,opt,name=stringStatResult,proto3" json:"stringStatResult,omitempty"`
}

func (x *FacetResult) Reset() {
	*x = FacetResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FacetResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FacetResult) ProtoMessage() {}

func (x *FacetResult) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FacetResult.ProtoReflect.Descriptor instead.
func (*FacetResult) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{8}
}

func (x *FacetResult) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *FacetResult) GetAxis() string {
	if x != nil {
		return x.Axis
	}
	return ""
}

func (x *FacetResult) GetBucketNo() uint32 {
	if x != nil {
		return x.BucketNo
	}
	return 0
}

func (x *FacetResult) GetBuckets() []*FacetBucket {
	if x != nil {
		return x.Buckets
	}
	return nil
}

func (x *FacetResult) GetStatName() string {
	if x != nil {
		return x.StatName
	}
	return ""
}

func (x *FacetResult) GetStringStatResult() string {
	if x != nil {
		return x.StringStatResult
	}
	return ""
}

type FieldHighlight struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FieldName string   `protobuf:"bytes,1,opt,name=fieldName,proto3" json:"fieldName,omitempty"`
	Snippets  []string `protobuf:"bytes,2,rep,name=snippets,proto3" json:"snippets,omitempty"`
	Language  string   `protobuf:"bytes,3,opt,name=language,proto3" json:"language,omitempty"`
}

func (x *FieldHighlight) Reset() {
	*x = FieldHighlight{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldHighlight) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldHighlight) ProtoMessage() {}

func (x *FieldHighlight) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldHighlight.ProtoReflect.Descriptor instead.
func (*FieldHighlight) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{9}
}

func (x *FieldHighlight) GetFieldName() string {
	if x != nil {
		return x.FieldName
	}
	return ""
}

func (x *FieldHighlight) GetSnippets() []string {
	if x != nil {
		return x.Snippets
	}
	return nil
}

func (x *FieldHighlight) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

type Highlight struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemId  string            `protobuf:"bytes,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	Matches []*FieldHighlight `protobuf:"bytes,2,rep,name=matches,proto3" json:"matches,omitempty"`
}

func (x *Highlight) Reset() {
	*x = Highlight{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Highlight) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Highlight) ProtoMessage() {}

func (x *Highlight) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Highlight.ProtoReflect.Descriptor instead.
func (*Highlight) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{10}
}

func (x *Highlight) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *Highlight) GetMatches() []*FieldHighlight {
	if x != nil {
		return x.Matches
	}
	return nil
}

type Diagnostics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParsingSucceeded bool     `protobuf:"varint,1,opt,name=parsing_succeeded,json=parsingSucceeded,proto3" json:"parsing_succeeded,omitempty"`
	ParsingErrors    []string `protobuf:"bytes,2,rep,name=parsing_errors,json=parsingErrors,proto3" json:"parsing_errors,omitempty"`
	CleanedQuery     string   `protobuf:"bytes,3,opt,name=cleaned_query,json=cleanedQuery,proto3" json:"cleaned_query,omitempty"`
	QueryWasCleaned  bool     `protobuf:"varint,4,opt,name=query_was_cleaned,json=queryWasCleaned,proto3" json:"query_was_cleaned,omitempty"`
	IgnoredErrors    []string `protobuf:"bytes,5,rep,name=ignored_errors,json=ignoredErrors,proto3" json:"ignored_errors,omitempty"`
}

func (x *Diagnostics) Reset() {
	*x = Diagnostics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_solr_solr_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Diagnostics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Diagnostics) ProtoMessage() {}

func (x *Diagnostics) ProtoReflect() protoreflect.Message {
	mi := &file_shared_solr_solr_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Diagnostics.ProtoReflect.Descriptor instead.
func (*Diagnostics) Descriptor() ([]byte, []int) {
	return file_shared_solr_solr_proto_rawDescGZIP(), []int{11}
}

func (x *Diagnostics) GetParsingSucceeded() bool {
	if x != nil {
		return x.ParsingSucceeded
	}
	return false
}

func (x *Diagnostics) GetParsingErrors() []string {
	if x != nil {
		return x.ParsingErrors
	}
	return nil
}

func (x *Diagnostics) GetCleanedQuery() string {
	if x != nil {
		return x.CleanedQuery
	}
	return ""
}

func (x *Diagnostics) GetQueryWasCleaned() bool {
	if x != nil {
		return x.QueryWasCleaned
	}
	return false
}

func (x *Diagnostics) GetIgnoredErrors() []string {
	if x != nil {
		return x.IgnoredErrors
	}
	return nil
}

var File_shared_solr_solr_proto protoreflect.FileDescriptor

var file_shared_solr_solr_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x73, 0x6f, 0x6c, 0x72, 0x2f, 0x73, 0x6f,
	0x6c, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6d, 0x65, 0x78, 0x2e, 0x76, 0x30,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x33, 0x0a, 0x07, 0x53,
	0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x78, 0x69, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x78, 0x69, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x22, 0x31, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x69,
	0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x61, 0x78, 0x22, 0xe3, 0x01, 0x0a, 0x0e, 0x41, 0x78, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x73,
	0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x78,
	0x69, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x78, 0x69, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65,
	0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x10, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x0d, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x65,
	0x78, 0x2e, 0x76, 0x30, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x52, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x29,
	0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x62, 0x69, 0x6e, 0x65, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6d, 0x62, 0x69, 0x6e,
	0x65, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x93, 0x01, 0x0a, 0x05, 0x46, 0x61,
	0x63, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x78, 0x69, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x78, 0x69, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x61,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74,
	0x61, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x74, 0x5f, 0x6f,
	0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x4f, 0x70, 0x22,
	0x66, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x22, 0x6d, 0x0a, 0x07, 0x44, 0x6f, 0x63, 0x49, 0x74,
	0x65, 0x6d, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x28, 0x0a, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d,
	0x65, 0x78, 0x2e, 0x76, 0x30, 0x2e, 0x44, 0x6f, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x62, 0x0a, 0x0d, 0x48, 0x69, 0x65, 0x72, 0x61, 0x72,
	0x63, 0x68, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69,
	0x73, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x69, 0x73,
	0x70, 0x6c, 0x61, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x22, 0x75, 0x0a, 0x0b, 0x46, 0x61,
	0x63, 0x65, 0x74, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3a, 0x0a, 0x0d, 0x68, 0x69, 0x65, 0x72, 0x61, 0x72, 0x63,
	0x68, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x52, 0x0d, 0x68, 0x69, 0x65, 0x72, 0x61, 0x72, 0x63, 0x68, 0x79, 0x49, 0x6e, 0x66,
	0x6f, 0x22, 0xc8, 0x01, 0x0a, 0x0b, 0x46, 0x61, 0x63, 0x65, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x78, 0x69, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x78, 0x69, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x4e, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x4e, 0x6f, 0x12, 0x2d, 0x0a, 0x07, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x76, 0x30, 0x2e,
	0x46, 0x61, 0x63, 0x65, 0x74, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x07, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x61, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x2a, 0x0a, 0x10, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x73, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x66, 0x0a, 0x0e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x48, 0x69, 0x67, 0x68, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x22, 0x55, 0x0a, 0x09, 0x48, 0x69, 0x67, 0x68, 0x6c, 0x69, 0x67, 0x68,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x07, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x78,
	0x2e, 0x76, 0x30, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x48, 0x69, 0x67, 0x68, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x52, 0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x22, 0xd9, 0x01, 0x0a, 0x0b,
	0x44, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x70,
	0x61, 0x72, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x70, 0x61, 0x72, 0x73, 0x69, 0x6e, 0x67, 0x53,
	0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x73,
	0x69, 0x6e, 0x67, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0d, 0x70, 0x61, 0x72, 0x73, 0x69, 0x6e, 0x67, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x12,
	0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x65, 0x61, 0x6e, 0x65, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x65, 0x61, 0x6e, 0x65, 0x64, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x2a, 0x0a, 0x11, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x77, 0x61,
	0x73, 0x5f, 0x63, 0x6c, 0x65, 0x61, 0x6e, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x57, 0x61, 0x73, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x65, 0x64,
	0x12, 0x25, 0x0a, 0x0e, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x64, 0x5f, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65,
	0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x73, 0x75, 0x6e, 0x64, 0x68, 0x65, 0x69, 0x74,
	0x73, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72, 0x6b, 0x69, 0x2d, 0x6d, 0x65, 0x78, 0x2d, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x78, 0x2f, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x2f, 0x73, 0x6f, 0x6c, 0x72, 0x3b, 0x73, 0x6f, 0x6c, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shared_solr_solr_proto_rawDescOnce sync.Once
	file_shared_solr_solr_proto_rawDescData = file_shared_solr_solr_proto_rawDesc
)

func file_shared_solr_solr_proto_rawDescGZIP() []byte {
	file_shared_solr_solr_proto_rawDescOnce.Do(func() {
		file_shared_solr_solr_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_solr_solr_proto_rawDescData)
	})
	return file_shared_solr_solr_proto_rawDescData
}

var file_shared_solr_solr_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_shared_solr_solr_proto_goTypes = []interface{}{
	(*Sorting)(nil),        // 0: mex.v0.Sorting
	(*StringRange)(nil),    // 1: mex.v0.StringRange
	(*AxisConstraint)(nil), // 2: mex.v0.AxisConstraint
	(*Facet)(nil),          // 3: mex.v0.Facet
	(*DocValue)(nil),       // 4: mex.v0.DocValue
	(*DocItem)(nil),        // 5: mex.v0.DocItem
	(*HierarchyInfo)(nil),  // 6: mex.v0.HierarchyInfo
	(*FacetBucket)(nil),    // 7: mex.v0.FacetBucket
	(*FacetResult)(nil),    // 8: mex.v0.FacetResult
	(*FieldHighlight)(nil), // 9: mex.v0.FieldHighlight
	(*Highlight)(nil),      // 10: mex.v0.Highlight
	(*Diagnostics)(nil),    // 11: mex.v0.Diagnostics
	(*anypb.Any)(nil),      // 12: google.protobuf.Any
}
var file_shared_solr_solr_proto_depIdxs = []int32{
	1,  // 0: mex.v0.AxisConstraint.string_ranges:type_name -> mex.v0.StringRange
	4,  // 1: mex.v0.DocItem.values:type_name -> mex.v0.DocValue
	12, // 2: mex.v0.FacetBucket.hierarchyInfo:type_name -> google.protobuf.Any
	7,  // 3: mex.v0.FacetResult.buckets:type_name -> mex.v0.FacetBucket
	9,  // 4: mex.v0.Highlight.matches:type_name -> mex.v0.FieldHighlight
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_shared_solr_solr_proto_init() }
func file_shared_solr_solr_proto_init() {
	if File_shared_solr_solr_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_solr_solr_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sorting); i {
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
		file_shared_solr_solr_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringRange); i {
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
		file_shared_solr_solr_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AxisConstraint); i {
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
		file_shared_solr_solr_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Facet); i {
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
		file_shared_solr_solr_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocValue); i {
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
		file_shared_solr_solr_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocItem); i {
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
		file_shared_solr_solr_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HierarchyInfo); i {
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
		file_shared_solr_solr_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FacetBucket); i {
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
		file_shared_solr_solr_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FacetResult); i {
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
		file_shared_solr_solr_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldHighlight); i {
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
		file_shared_solr_solr_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Highlight); i {
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
		file_shared_solr_solr_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Diagnostics); i {
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
			RawDescriptor: file_shared_solr_solr_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_solr_solr_proto_goTypes,
		DependencyIndexes: file_shared_solr_solr_proto_depIdxs,
		MessageInfos:      file_shared_solr_solr_proto_msgTypes,
	}.Build()
	File_shared_solr_solr_proto = out.File
	file_shared_solr_solr_proto_rawDesc = nil
	file_shared_solr_solr_proto_goTypes = nil
	file_shared_solr_solr_proto_depIdxs = nil
}
