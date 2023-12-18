package authn

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/known/securitypb"
)

type APIKeyAuthenticator struct {
	tenantID string
	keyRoles *auth.ApiKeys
}

func NewAPIKeyAuthenticator(tenantID string, keysPrivilgesBytes []byte) (*APIKeyAuthenticator, error) {
	if len(keysPrivilgesBytes) == 0 {
		return &APIKeyAuthenticator{
			tenantID: tenantID,
			keyRoles: &auth.ApiKeys{},
		}, nil
	}

	var keyPrivileges auth.ApiKeys
	err := protojson.Unmarshal(keysPrivilgesBytes, &keyPrivileges)
	if err != nil {
		return nil, err
	}

	return &APIKeyAuthenticator{
		tenantID: tenantID,
		keyRoles: &keyPrivileges,
	}, nil
}

func (a *APIKeyAuthenticator) Authenticate(ctx context.Context, request any) (*securitypb.UserWithRoles, error) {
	authHeader, err := extractAuthorizationHeader(ctx)
	if err != nil {
		return nil, err
	}

	if authHeader.tokenType != "apikey" {
		return nil, status.Error(codes.Unauthenticated, "authorization header does not contain an API key")
	}

	roleName, ok := a.keyRoles.KeysRoles[authHeader.tokenValue]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unknown API key")
	}

	user := securitypb.UserWithRoles{
		TenantId: a.tenantID,
		Roles:    []string{roleName},
	}

	return &user, nil
}

func (a *APIKeyAuthenticator) Count() int {
	if a.keyRoles == nil {
		return 0
	}

	if a.keyRoles.KeysRoles == nil {
		return 0
	}

	return len(a.keyRoles.KeysRoles)
}
