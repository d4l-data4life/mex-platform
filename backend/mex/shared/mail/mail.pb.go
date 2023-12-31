// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.0
// source: shared/mail/mail.proto

package mail

import (
	_ "github.com/d4l-data4life/mex/mex/shared/known/securitypb"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RecipientType int32

const (
	RecipientType_TO  RecipientType = 0
	RecipientType_CC  RecipientType = 1
	RecipientType_BCC RecipientType = 2
)

// Enum value maps for RecipientType.
var (
	RecipientType_name = map[int32]string{
		0: "TO",
		1: "CC",
		2: "BCC",
	}
	RecipientType_value = map[string]int32{
		"TO":  0,
		"CC":  1,
		"BCC": 2,
	}
)

func (x RecipientType) Enum() *RecipientType {
	p := new(RecipientType)
	*p = x
	return p
}

func (x RecipientType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RecipientType) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_mail_mail_proto_enumTypes[0].Descriptor()
}

func (RecipientType) Type() protoreflect.EnumType {
	return &file_shared_mail_mail_proto_enumTypes[0]
}

func (x RecipientType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RecipientType.Descriptor instead.
func (RecipientType) EnumDescriptor() ([]byte, []int) {
	return file_shared_mail_mail_proto_rawDescGZIP(), []int{0}
}

type StaticContact struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *StaticContact) Reset() {
	*x = StaticContact{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_mail_mail_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StaticContact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaticContact) ProtoMessage() {}

func (x *StaticContact) ProtoReflect() protoreflect.Message {
	mi := &file_shared_mail_mail_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaticContact.ProtoReflect.Descriptor instead.
func (*StaticContact) Descriptor() ([]byte, []int) {
	return file_shared_mail_mail_proto_rawDescGZIP(), []int{0}
}

func (x *StaticContact) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *StaticContact) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type FieldNamesContact struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NameField  string `protobuf:"bytes,1,opt,name=name_field,json=nameField,proto3" json:"name_field,omitempty"`
	EmailField string `protobuf:"bytes,2,opt,name=email_field,json=emailField,proto3" json:"email_field,omitempty"`
}

func (x *FieldNamesContact) Reset() {
	*x = FieldNamesContact{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_mail_mail_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldNamesContact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldNamesContact) ProtoMessage() {}

func (x *FieldNamesContact) ProtoReflect() protoreflect.Message {
	mi := &file_shared_mail_mail_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldNamesContact.ProtoReflect.Descriptor instead.
func (*FieldNamesContact) Descriptor() ([]byte, []int) {
	return file_shared_mail_mail_proto_rawDescGZIP(), []int{1}
}

func (x *FieldNamesContact) GetNameField() string {
	if x != nil {
		return x.NameField
	}
	return ""
}

func (x *FieldNamesContact) GetEmailField() string {
	if x != nil {
		return x.EmailField
	}
	return ""
}

type Sender struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to ContactType:
	//
	//	*Sender_Static
	//	*Sender_FormData
	ContactType isSender_ContactType `protobuf_oneof:"contact_type"`
}

func (x *Sender) Reset() {
	*x = Sender{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_mail_mail_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sender) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sender) ProtoMessage() {}

func (x *Sender) ProtoReflect() protoreflect.Message {
	mi := &file_shared_mail_mail_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sender.ProtoReflect.Descriptor instead.
func (*Sender) Descriptor() ([]byte, []int) {
	return file_shared_mail_mail_proto_rawDescGZIP(), []int{2}
}

func (m *Sender) GetContactType() isSender_ContactType {
	if m != nil {
		return m.ContactType
	}
	return nil
}

func (x *Sender) GetStatic() *StaticContact {
	if x, ok := x.GetContactType().(*Sender_Static); ok {
		return x.Static
	}
	return nil
}

func (x *Sender) GetFormData() *FieldNamesContact {
	if x, ok := x.GetContactType().(*Sender_FormData); ok {
		return x.FormData
	}
	return nil
}

type isSender_ContactType interface {
	isSender_ContactType()
}

type Sender_Static struct {
	Static *StaticContact `protobuf:"bytes,1,opt,name=static,proto3,oneof"`
}

type Sender_FormData struct {
	FormData *FieldNamesContact `protobuf:"bytes,2,opt,name=form_data,json=formData,proto3,oneof"`
}

func (*Sender_Static) isSender_ContactType() {}

func (*Sender_FormData) isSender_ContactType() {}

type Recipient struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type RecipientType `protobuf:"varint,1,opt,name=type,proto3,enum=d4l.mex.mail.RecipientType" json:"type,omitempty"`
	// Types that are assignable to ContactType:
	//
	//	*Recipient_Static
	//	*Recipient_FormData
	//	*Recipient_ContextItem
	//	*Recipient_ContactItem
	ContactType isRecipient_ContactType `protobuf_oneof:"contact_type"`
}

func (x *Recipient) Reset() {
	*x = Recipient{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_mail_mail_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Recipient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Recipient) ProtoMessage() {}

func (x *Recipient) ProtoReflect() protoreflect.Message {
	mi := &file_shared_mail_mail_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Recipient.ProtoReflect.Descriptor instead.
func (*Recipient) Descriptor() ([]byte, []int) {
	return file_shared_mail_mail_proto_rawDescGZIP(), []int{3}
}

func (x *Recipient) GetType() RecipientType {
	if x != nil {
		return x.Type
	}
	return RecipientType_TO
}

func (m *Recipient) GetContactType() isRecipient_ContactType {
	if m != nil {
		return m.ContactType
	}
	return nil
}

func (x *Recipient) GetStatic() *StaticContact {
	if x, ok := x.GetContactType().(*Recipient_Static); ok {
		return x.Static
	}
	return nil
}

func (x *Recipient) GetFormData() *FieldNamesContact {
	if x, ok := x.GetContactType().(*Recipient_FormData); ok {
		return x.FormData
	}
	return nil
}

func (x *Recipient) GetContextItem() *FieldNamesContact {
	if x, ok := x.GetContactType().(*Recipient_ContextItem); ok {
		return x.ContextItem
	}
	return nil
}

func (x *Recipient) GetContactItem() *FieldNamesContact {
	if x, ok := x.GetContactType().(*Recipient_ContactItem); ok {
		return x.ContactItem
	}
	return nil
}

type isRecipient_ContactType interface {
	isRecipient_ContactType()
}

type Recipient_Static struct {
	Static *StaticContact `protobuf:"bytes,2,opt,name=static,proto3,oneof"`
}

type Recipient_FormData struct {
	FormData *FieldNamesContact `protobuf:"bytes,3,opt,name=form_data,json=formData,proto3,oneof"`
}

type Recipient_ContextItem struct {
	ContextItem *FieldNamesContact `protobuf:"bytes,4,opt,name=context_item,json=contextItem,proto3,oneof"`
}

type Recipient_ContactItem struct {
	ContactItem *FieldNamesContact `protobuf:"bytes,5,opt,name=contact_item,json=contactItem,proto3,oneof"`
}

func (*Recipient_Static) isRecipient_ContactType() {}

func (*Recipient_FormData) isRecipient_ContactType() {}

func (*Recipient_ContextItem) isRecipient_ContactType() {}

func (*Recipient_ContactItem) isRecipient_ContactType() {}

type MailTemplate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Sender         *Sender      `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Recipients     []*Recipient `protobuf:"bytes,3,rep,name=recipients,proto3" json:"recipients,omitempty"`
	TemplateEngine string       `protobuf:"bytes,4,opt,name=template_engine,json=templateEngine,proto3" json:"template_engine,omitempty"`
	Subject        string       `protobuf:"bytes,5,opt,name=subject,proto3" json:"subject,omitempty"`
	TextBody       string       `protobuf:"bytes,6,opt,name=text_body,json=textBody,proto3" json:"text_body,omitempty"`
	HtmlBody       string       `protobuf:"bytes,7,opt,name=html_body,json=htmlBody,proto3" json:"html_body,omitempty"`
}

func (x *MailTemplate) Reset() {
	*x = MailTemplate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_mail_mail_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailTemplate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailTemplate) ProtoMessage() {}

func (x *MailTemplate) ProtoReflect() protoreflect.Message {
	mi := &file_shared_mail_mail_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailTemplate.ProtoReflect.Descriptor instead.
func (*MailTemplate) Descriptor() ([]byte, []int) {
	return file_shared_mail_mail_proto_rawDescGZIP(), []int{4}
}

func (x *MailTemplate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MailTemplate) GetSender() *Sender {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *MailTemplate) GetRecipients() []*Recipient {
	if x != nil {
		return x.Recipients
	}
	return nil
}

func (x *MailTemplate) GetTemplateEngine() string {
	if x != nil {
		return x.TemplateEngine
	}
	return ""
}

func (x *MailTemplate) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *MailTemplate) GetTextBody() string {
	if x != nil {
		return x.TextBody
	}
	return ""
}

func (x *MailTemplate) GetHtmlBody() string {
	if x != nil {
		return x.HtmlBody
	}
	return ""
}

type MailOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId        string           `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Subject        string           `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	TextBody       string           `protobuf:"bytes,3,opt,name=text_body,json=textBody,proto3" json:"text_body,omitempty"`
	HtmlBody       string           `protobuf:"bytes,4,opt,name=html_body,json=htmlBody,proto3" json:"html_body,omitempty"`
	Sender         *StaticContact   `protobuf:"bytes,5,opt,name=sender,proto3" json:"sender,omitempty"`
	RecipientsTo   []*StaticContact `protobuf:"bytes,6,rep,name=recipients_to,json=recipientsTo,proto3" json:"recipients_to,omitempty"`
	RecipientsCc   []*StaticContact `protobuf:"bytes,7,rep,name=recipients_cc,json=recipientsCc,proto3" json:"recipients_cc,omitempty"`
	RecipientsBcc  []*StaticContact `protobuf:"bytes,8,rep,name=recipients_bcc,json=recipientsBcc,proto3" json:"recipients_bcc,omitempty"`
	TemplateEngine string           `protobuf:"bytes,9,opt,name=template_engine,json=templateEngine,proto3" json:"template_engine,omitempty"`
}

func (x *MailOrder) Reset() {
	*x = MailOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_mail_mail_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailOrder) ProtoMessage() {}

func (x *MailOrder) ProtoReflect() protoreflect.Message {
	mi := &file_shared_mail_mail_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailOrder.ProtoReflect.Descriptor instead.
func (*MailOrder) Descriptor() ([]byte, []int) {
	return file_shared_mail_mail_proto_rawDescGZIP(), []int{5}
}

func (x *MailOrder) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *MailOrder) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *MailOrder) GetTextBody() string {
	if x != nil {
		return x.TextBody
	}
	return ""
}

func (x *MailOrder) GetHtmlBody() string {
	if x != nil {
		return x.HtmlBody
	}
	return ""
}

func (x *MailOrder) GetSender() *StaticContact {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *MailOrder) GetRecipientsTo() []*StaticContact {
	if x != nil {
		return x.RecipientsTo
	}
	return nil
}

func (x *MailOrder) GetRecipientsCc() []*StaticContact {
	if x != nil {
		return x.RecipientsCc
	}
	return nil
}

func (x *MailOrder) GetRecipientsBcc() []*StaticContact {
	if x != nil {
		return x.RecipientsBcc
	}
	return nil
}

func (x *MailOrder) GetTemplateEngine() string {
	if x != nil {
		return x.TemplateEngine
	}
	return ""
}

var File_shared_mail_mail_proto protoreflect.FileDescriptor

var file_shared_mail_mail_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x2f, 0x6d, 0x61,
	0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65,
	0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x12, 0x64, 0x34, 0x6c, 0x2f, 0x73, 0x65, 0x63, 0x75,
	0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x39, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x53,
	0x0a, 0x11, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x22, 0x8f, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x35,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x48, 0x00, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x12, 0x3e, 0x0a, 0x09, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d,
	0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x48, 0x00, 0x52, 0x08, 0x66, 0x6f, 0x72,
	0x6d, 0x44, 0x61, 0x74, 0x61, 0x42, 0x0e, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0xcf, 0x02, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1b, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c,
	0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d,
	0x61, 0x69, 0x6c, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x12, 0x3e, 0x0a, 0x09, 0x66,
	0x6f, 0x72, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x48,
	0x00, 0x52, 0x08, 0x66, 0x6f, 0x72, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x12, 0x44, 0x0a, 0x0c, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x61,
	0x63, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x44, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x74, 0x65,
	0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65,
	0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x42, 0x0e, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x61,
	0x63, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0x86, 0x02, 0x0a, 0x0c, 0x4d, 0x61, 0x69, 0x6c,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x06,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64,
	0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x0a, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x52, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x74, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x78, 0x74, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x74, 0x6d, 0x6c, 0x5f, 0x62, 0x6f, 0x64, 0x79,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x74, 0x6d, 0x6c, 0x42, 0x6f, 0x64, 0x79,
	0x22, 0xa0, 0x03, 0x0a, 0x09, 0x4d, 0x61, 0x69, 0x6c, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x19,
	0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x62, 0x6f, 0x64, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x78, 0x74, 0x42, 0x6f, 0x64, 0x79,
	0x12, 0x1b, 0x0a, 0x09, 0x68, 0x74, 0x6d, 0x6c, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x74, 0x6d, 0x6c, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x33, 0x0a,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x12, 0x40, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73,
	0x5f, 0x74, 0x6f, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x34, 0x6c, 0x2e,
	0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x0c, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x73, 0x54, 0x6f, 0x12, 0x40, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x73, 0x5f, 0x63, 0x63, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x34,
	0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69,
	0x63, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x0c, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x73, 0x43, 0x63, 0x12, 0x42, 0x0a, 0x0e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x73, 0x5f, 0x62, 0x63, 0x63, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x0d, 0x72, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x63, 0x63, 0x12, 0x27, 0x0a, 0x0f, 0x74, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2a, 0x28, 0x0a, 0x0d, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x54, 0x4f, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02,
	0x43, 0x43, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x43, 0x43, 0x10, 0x02, 0x42, 0x43, 0x5a,
	0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x73, 0x75,
	0x6e, 0x64, 0x68, 0x65, 0x69, 0x74, 0x73, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72, 0x6b, 0x69,
	0x2d, 0x6d, 0x65, 0x78, 0x2d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65,
	0x78, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x3b, 0x6d, 0x61,
	0x69, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shared_mail_mail_proto_rawDescOnce sync.Once
	file_shared_mail_mail_proto_rawDescData = file_shared_mail_mail_proto_rawDesc
)

func file_shared_mail_mail_proto_rawDescGZIP() []byte {
	file_shared_mail_mail_proto_rawDescOnce.Do(func() {
		file_shared_mail_mail_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_mail_mail_proto_rawDescData)
	})
	return file_shared_mail_mail_proto_rawDescData
}

var file_shared_mail_mail_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_shared_mail_mail_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_shared_mail_mail_proto_goTypes = []interface{}{
	(RecipientType)(0),        // 0: d4l.mex.mail.RecipientType
	(*StaticContact)(nil),     // 1: d4l.mex.mail.StaticContact
	(*FieldNamesContact)(nil), // 2: d4l.mex.mail.FieldNamesContact
	(*Sender)(nil),            // 3: d4l.mex.mail.Sender
	(*Recipient)(nil),         // 4: d4l.mex.mail.Recipient
	(*MailTemplate)(nil),      // 5: d4l.mex.mail.MailTemplate
	(*MailOrder)(nil),         // 6: d4l.mex.mail.MailOrder
}
var file_shared_mail_mail_proto_depIdxs = []int32{
	1,  // 0: d4l.mex.mail.Sender.static:type_name -> d4l.mex.mail.StaticContact
	2,  // 1: d4l.mex.mail.Sender.form_data:type_name -> d4l.mex.mail.FieldNamesContact
	0,  // 2: d4l.mex.mail.Recipient.type:type_name -> d4l.mex.mail.RecipientType
	1,  // 3: d4l.mex.mail.Recipient.static:type_name -> d4l.mex.mail.StaticContact
	2,  // 4: d4l.mex.mail.Recipient.form_data:type_name -> d4l.mex.mail.FieldNamesContact
	2,  // 5: d4l.mex.mail.Recipient.context_item:type_name -> d4l.mex.mail.FieldNamesContact
	2,  // 6: d4l.mex.mail.Recipient.contact_item:type_name -> d4l.mex.mail.FieldNamesContact
	3,  // 7: d4l.mex.mail.MailTemplate.sender:type_name -> d4l.mex.mail.Sender
	4,  // 8: d4l.mex.mail.MailTemplate.recipients:type_name -> d4l.mex.mail.Recipient
	1,  // 9: d4l.mex.mail.MailOrder.sender:type_name -> d4l.mex.mail.StaticContact
	1,  // 10: d4l.mex.mail.MailOrder.recipients_to:type_name -> d4l.mex.mail.StaticContact
	1,  // 11: d4l.mex.mail.MailOrder.recipients_cc:type_name -> d4l.mex.mail.StaticContact
	1,  // 12: d4l.mex.mail.MailOrder.recipients_bcc:type_name -> d4l.mex.mail.StaticContact
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_shared_mail_mail_proto_init() }
func file_shared_mail_mail_proto_init() {
	if File_shared_mail_mail_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_mail_mail_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StaticContact); i {
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
		file_shared_mail_mail_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldNamesContact); i {
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
		file_shared_mail_mail_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sender); i {
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
		file_shared_mail_mail_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Recipient); i {
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
		file_shared_mail_mail_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailTemplate); i {
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
		file_shared_mail_mail_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailOrder); i {
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
	file_shared_mail_mail_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Sender_Static)(nil),
		(*Sender_FormData)(nil),
	}
	file_shared_mail_mail_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Recipient_Static)(nil),
		(*Recipient_FormData)(nil),
		(*Recipient_ContextItem)(nil),
		(*Recipient_ContactItem)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shared_mail_mail_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_mail_mail_proto_goTypes,
		DependencyIndexes: file_shared_mail_mail_proto_depIdxs,
		EnumInfos:         file_shared_mail_mail_proto_enumTypes,
		MessageInfos:      file_shared_mail_mail_proto_msgTypes,
	}.Build()
	File_shared_mail_mail_proto = out.File
	file_shared_mail_mail_proto_rawDesc = nil
	file_shared_mail_mail_proto_goTypes = nil
	file_shared_mail_mail_proto_depIdxs = nil
}
