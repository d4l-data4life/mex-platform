package authn

import (
	"context"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/known/securitypb"
)

type noneAuthenticator struct {
	tenantID string
}

func NewNoneAuthenticator(tenantID string) auth.RequestAuthenticator {
	return &noneAuthenticator{tenantID: tenantID}
}

func (a *noneAuthenticator) Authenticate(ctx context.Context, request any) (*securitypb.UserWithRoles, error) {
	return &securitypb.UserWithRoles{
		TenantId: a.tenantID,
		UserId:   "anonymous",
		Roles:    []string{},
	}, nil
}
