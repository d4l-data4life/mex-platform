// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.0
// source: services/auth/endpoints/auth/auth.proto

package pbAuth

import (
	_ "github.com/d4l-data4life/mex/mex/shared/known/securitypb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type AuthorizeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResponseType        string `protobuf:"bytes,1,opt,name=response_type,json=responseType,proto3" json:"response_type,omitempty"`
	ClientId            string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	RedirectUri         string `protobuf:"bytes,3,opt,name=redirect_uri,json=redirectUri,proto3" json:"redirect_uri,omitempty"`
	Scope               string `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	CodeChallenge       string `protobuf:"bytes,5,opt,name=code_challenge,json=codeChallenge,proto3" json:"code_challenge,omitempty"`
	CodeChallengeMethod string `protobuf:"bytes,6,opt,name=code_challenge_method,json=codeChallengeMethod,proto3" json:"code_challenge_method,omitempty"`
	State               string `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
	ResponseMode        string `protobuf:"bytes,9,opt,name=response_mode,json=responseMode,proto3" json:"response_mode,omitempty"`
}

func (x *AuthorizeRequest) Reset() {
	*x = AuthorizeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorizeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizeRequest) ProtoMessage() {}

func (x *AuthorizeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizeRequest.ProtoReflect.Descriptor instead.
func (*AuthorizeRequest) Descriptor() ([]byte, []int) {
	return file_services_auth_endpoints_auth_auth_proto_rawDescGZIP(), []int{0}
}

func (x *AuthorizeRequest) GetResponseType() string {
	if x != nil {
		return x.ResponseType
	}
	return ""
}

func (x *AuthorizeRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *AuthorizeRequest) GetRedirectUri() string {
	if x != nil {
		return x.RedirectUri
	}
	return ""
}

func (x *AuthorizeRequest) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *AuthorizeRequest) GetCodeChallenge() string {
	if x != nil {
		return x.CodeChallenge
	}
	return ""
}

func (x *AuthorizeRequest) GetCodeChallengeMethod() string {
	if x != nil {
		return x.CodeChallengeMethod
	}
	return ""
}

func (x *AuthorizeRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *AuthorizeRequest) GetResponseMode() string {
	if x != nil {
		return x.ResponseMode
	}
	return ""
}

type AuthorizeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AuthorizeResponse) Reset() {
	*x = AuthorizeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorizeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizeResponse) ProtoMessage() {}

func (x *AuthorizeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizeResponse.ProtoReflect.Descriptor instead.
func (*AuthorizeResponse) Descriptor() ([]byte, []int) {
	return file_services_auth_endpoints_auth_auth_proto_rawDescGZIP(), []int{1}
}

type TokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	RedirectUri  string `protobuf:"bytes,2,opt,name=redirect_uri,json=redirectUri,proto3" json:"redirect_uri,omitempty"`
	GrantType    string `protobuf:"bytes,3,opt,name=grant_type,json=grantType,proto3" json:"grant_type,omitempty"`
	CodeVerifier string `protobuf:"bytes,4,opt,name=code_verifier,json=codeVerifier,proto3" json:"code_verifier,omitempty"`
	Code         string `protobuf:"bytes,5,opt,name=code,proto3" json:"code,omitempty"`
	Scope        string `protobuf:"bytes,6,opt,name=scope,proto3" json:"scope,omitempty"`
	State        string `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
	ClientSecret string `protobuf:"bytes,9,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
	RefreshToken string `protobuf:"bytes,10,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
}

func (x *TokenRequest) Reset() {
	*x = TokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenRequest) ProtoMessage() {}

func (x *TokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenRequest.ProtoReflect.Descriptor instead.
func (*TokenRequest) Descriptor() ([]byte, []int) {
	return file_services_auth_endpoints_auth_auth_proto_rawDescGZIP(), []int{2}
}

func (x *TokenRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *TokenRequest) GetRedirectUri() string {
	if x != nil {
		return x.RedirectUri
	}
	return ""
}

func (x *TokenRequest) GetGrantType() string {
	if x != nil {
		return x.GrantType
	}
	return ""
}

func (x *TokenRequest) GetCodeVerifier() string {
	if x != nil {
		return x.CodeVerifier
	}
	return ""
}

func (x *TokenRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *TokenRequest) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *TokenRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *TokenRequest) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *TokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

type TokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken  string `protobuf:"bytes,1,opt,name=access_token,proto3" json:"access_token,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token,proto3" json:"refresh_token,omitempty"`
	TokenType    string `protobuf:"bytes,3,opt,name=token_type,proto3" json:"token_type,omitempty"`
	ExpiresIn    uint32 `protobuf:"varint,4,opt,name=expires_in,proto3" json:"expires_in,omitempty"`
}

func (x *TokenResponse) Reset() {
	*x = TokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenResponse) ProtoMessage() {}

func (x *TokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenResponse.ProtoReflect.Descriptor instead.
func (*TokenResponse) Descriptor() ([]byte, []int) {
	return file_services_auth_endpoints_auth_auth_proto_rawDescGZIP(), []int{3}
}

func (x *TokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *TokenResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *TokenResponse) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *TokenResponse) GetExpiresIn() uint32 {
	if x != nil {
		return x.ExpiresIn
	}
	return 0
}

type KeysRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *KeysRequest) Reset() {
	*x = KeysRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeysRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeysRequest) ProtoMessage() {}

func (x *KeysRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeysRequest.ProtoReflect.Descriptor instead.
func (*KeysRequest) Descriptor() ([]byte, []int) {
	return file_services_auth_endpoints_auth_auth_proto_rawDescGZIP(), []int{4}
}

type KeysResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys []*KeysResponse_Key `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *KeysResponse) Reset() {
	*x = KeysResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeysResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeysResponse) ProtoMessage() {}

func (x *KeysResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeysResponse.ProtoReflect.Descriptor instead.
func (*KeysResponse) Descriptor() ([]byte, []int) {
	return file_services_auth_endpoints_auth_auth_proto_rawDescGZIP(), []int{5}
}

func (x *KeysResponse) GetKeys() []*KeysResponse_Key {
	if x != nil {
		return x.Keys
	}
	return nil
}

type KeysResponse_Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kty string `protobuf:"bytes,1,opt,name=kty,proto3" json:"kty,omitempty"`
	Kid string `protobuf:"bytes,2,opt,name=kid,proto3" json:"kid,omitempty"`
	Use string `protobuf:"bytes,3,opt,name=use,proto3" json:"use,omitempty"`
	Alg string `protobuf:"bytes,4,opt,name=alg,proto3" json:"alg,omitempty"`
	E   string `protobuf:"bytes,5,opt,name=e,proto3" json:"e,omitempty"`
	N   string `protobuf:"bytes,6,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *KeysResponse_Key) Reset() {
	*x = KeysResponse_Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeysResponse_Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeysResponse_Key) ProtoMessage() {}

func (x *KeysResponse_Key) ProtoReflect() protoreflect.Message {
	mi := &file_services_auth_endpoints_auth_auth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeysResponse_Key.ProtoReflect.Descriptor instead.
func (*KeysResponse_Key) Descriptor() ([]byte, []int) {
	return file_services_auth_endpoints_auth_auth_proto_rawDescGZIP(), []int{5, 0}
}

func (x *KeysResponse_Key) GetKty() string {
	if x != nil {
		return x.Kty
	}
	return ""
}

func (x *KeysResponse_Key) GetKid() string {
	if x != nil {
		return x.Kid
	}
	return ""
}

func (x *KeysResponse_Key) GetUse() string {
	if x != nil {
		return x.Use
	}
	return ""
}

func (x *KeysResponse_Key) GetAlg() string {
	if x != nil {
		return x.Alg
	}
	return ""
}

func (x *KeysResponse_Key) GetE() string {
	if x != nil {
		return x.E
	}
	return ""
}

func (x *KeysResponse_Key) GetN() string {
	if x != nil {
		return x.N
	}
	return ""
}

var File_services_auth_endpoints_auth_auth_proto protoreflect.FileDescriptor

var file_services_auth_endpoints_auth_auth_proto_rawDesc = []byte{
	0x0a, 0x27, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x64, 0x34, 0x6c, 0x2e, 0x6d,
	0x65, 0x78, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x1a, 0x12, 0x64, 0x34, 0x6c, 0x2f, 0x73, 0x65, 0x63,
	0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa3, 0x02, 0x0a, 0x10, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23,
	0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x69,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x55, 0x72, 0x69, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x64,
	0x65, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x63, 0x6f, 0x64, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x12, 0x32, 0x0a, 0x15, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x13, 0x63, 0x6f, 0x64, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x22,
	0x13, 0x0a, 0x11, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x9c, 0x02, 0x0a, 0x0c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75,
	0x72, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x55, 0x72, 0x69, 0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x61, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x76, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x64,
	0x65, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x23,
	0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x99, 0x01, 0x0a, 0x0d, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x22,
	0x0d, 0x0a, 0x0b, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xad,
	0x01, 0x0a, 0x0c, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x32, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x4b, 0x65, 0x79,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x04, 0x6b,
	0x65, 0x79, 0x73, 0x1a, 0x69, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x74,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x74, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x69, 0x64, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x73, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x61, 0x6c, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61,
	0x6c, 0x67, 0x12, 0x0c, 0x0a, 0x01, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x65,
	0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x6e, 0x32, 0xbe,
	0x02, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x71, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x23, 0x98, 0xf1, 0x04, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x19, 0x12, 0x17, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x30, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68,
	0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x12, 0x64, 0x0a, 0x05, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x22, 0x98, 0xf1,
	0x04, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x22, 0x13, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x30, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x5d, 0x0a, 0x04, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x19, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d,
	0x65, 0x78, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x64, 0x34, 0x6c, 0x2e, 0x6d, 0x65, 0x78, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x1e, 0x98, 0xf1, 0x04, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x12, 0x12, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x30, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x6b, 0x65, 0x79, 0x73, 0x42,
	0x59, 0x5a, 0x57, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65,
	0x73, 0x75, 0x6e, 0x64, 0x68, 0x65, 0x69, 0x74, 0x73, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72,
	0x6b, 0x69, 0x2d, 0x6d, 0x65, 0x78, 0x2d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x6d, 0x65, 0x78, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x2f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68,
	0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x41, 0x75, 0x74, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_services_auth_endpoints_auth_auth_proto_rawDescOnce sync.Once
	file_services_auth_endpoints_auth_auth_proto_rawDescData = file_services_auth_endpoints_auth_auth_proto_rawDesc
)

func file_services_auth_endpoints_auth_auth_proto_rawDescGZIP() []byte {
	file_services_auth_endpoints_auth_auth_proto_rawDescOnce.Do(func() {
		file_services_auth_endpoints_auth_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_auth_endpoints_auth_auth_proto_rawDescData)
	})
	return file_services_auth_endpoints_auth_auth_proto_rawDescData
}

var file_services_auth_endpoints_auth_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_services_auth_endpoints_auth_auth_proto_goTypes = []interface{}{
	(*AuthorizeRequest)(nil),  // 0: d4l.mex.auth.AuthorizeRequest
	(*AuthorizeResponse)(nil), // 1: d4l.mex.auth.AuthorizeResponse
	(*TokenRequest)(nil),      // 2: d4l.mex.auth.TokenRequest
	(*TokenResponse)(nil),     // 3: d4l.mex.auth.TokenResponse
	(*KeysRequest)(nil),       // 4: d4l.mex.auth.KeysRequest
	(*KeysResponse)(nil),      // 5: d4l.mex.auth.KeysResponse
	(*KeysResponse_Key)(nil),  // 6: d4l.mex.auth.KeysResponse.Key
}
var file_services_auth_endpoints_auth_auth_proto_depIdxs = []int32{
	6, // 0: d4l.mex.auth.KeysResponse.keys:type_name -> d4l.mex.auth.KeysResponse.Key
	0, // 1: d4l.mex.auth.Auth.Authorize:input_type -> d4l.mex.auth.AuthorizeRequest
	2, // 2: d4l.mex.auth.Auth.Token:input_type -> d4l.mex.auth.TokenRequest
	4, // 3: d4l.mex.auth.Auth.Keys:input_type -> d4l.mex.auth.KeysRequest
	1, // 4: d4l.mex.auth.Auth.Authorize:output_type -> d4l.mex.auth.AuthorizeResponse
	3, // 5: d4l.mex.auth.Auth.Token:output_type -> d4l.mex.auth.TokenResponse
	5, // 6: d4l.mex.auth.Auth.Keys:output_type -> d4l.mex.auth.KeysResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_services_auth_endpoints_auth_auth_proto_init() }
func file_services_auth_endpoints_auth_auth_proto_init() {
	if File_services_auth_endpoints_auth_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_auth_endpoints_auth_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorizeRequest); i {
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
		file_services_auth_endpoints_auth_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorizeResponse); i {
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
		file_services_auth_endpoints_auth_auth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenRequest); i {
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
		file_services_auth_endpoints_auth_auth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenResponse); i {
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
		file_services_auth_endpoints_auth_auth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeysRequest); i {
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
		file_services_auth_endpoints_auth_auth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeysResponse); i {
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
		file_services_auth_endpoints_auth_auth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeysResponse_Key); i {
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
			RawDescriptor: file_services_auth_endpoints_auth_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_auth_endpoints_auth_auth_proto_goTypes,
		DependencyIndexes: file_services_auth_endpoints_auth_auth_proto_depIdxs,
		MessageInfos:      file_services_auth_endpoints_auth_auth_proto_msgTypes,
	}.Build()
	File_services_auth_endpoints_auth_auth_proto = out.File
	file_services_auth_endpoints_auth_auth_proto_rawDesc = nil
	file_services_auth_endpoints_auth_auth_proto_goTypes = nil
	file_services_auth_endpoints_auth_auth_proto_depIdxs = nil
}
