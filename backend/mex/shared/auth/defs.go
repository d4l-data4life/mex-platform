package auth

import (
	"context"

	"github.com/d4l-data4life/mex/mex/shared/known/securitypb"
)

type RequestAuthenticator interface {
	Authenticate(ctx context.Context, request any) (*securitypb.UserWithRoles, error)
}

type RequestAuthenticatorRegistry = map[securitypb.AuthenticationType]RequestAuthenticator

type keyType string

const ContextKeyMexUser keyType = "mex-user"
