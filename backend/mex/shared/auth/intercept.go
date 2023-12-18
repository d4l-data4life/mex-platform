package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/known/securitypb"
)

func NewInterceptor(registry RequestAuthenticatorRegistry, privMgr *PrivMgr) grpc.UnaryServerInterceptor {
	if registry == nil {
		panic("registry is nil")
	}

	if privMgr == nil {
		panic("privilege manager is nil")
	}

	// Cache the methods' privilege masks
	mu := sync.Mutex{}
	methodPrivileges := map[string]PrivMask{}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		authnType, err := getMethodAnnotation[securitypb.AuthenticationType](info.FullMethod, securitypb.E_AuthnType)
		if err != nil {
			panic("no security annotation: " + info.FullMethod)
		}

		requestAuthenticator, ok := registry[*authnType]
		if !ok {
			return nil, fmt.Errorf("unknown authentication type: %v", *authnType)
		}

		userWithRoles, err := requestAuthenticator.Authenticate(ctx, req)
		if err != nil {
			return nil, err
		}

		// At this point we are properly authenticated.
		// Now let's authorize, that is, check the required vs actual privileges.

		userPrivilegesMask, err := privMgr.ResolveRoles(userWithRoles.Roles)
		if err != nil {
			return nil, err
		}

		// Prevent concurrent modifications to the `methodPrivileges` map below when
		// handling concurrent requests.
		mu.Lock()
		defer mu.Unlock()

		requiredPrivilegesMask, ok := methodPrivileges[info.FullMethod]
		if !ok {
			requiredPrivileges, err := getMethodAnnotation[[]*securitypb.Privilege](info.FullMethod, securitypb.E_RequiredPrivileges)
			if err != nil {
				requiredPrivileges = &[]*securitypb.Privilege{}
			}

			for _, requiredPrivilege := range *requiredPrivileges {
				requiredPrivilegesMask |= privMgr.MustPrivMask(requiredPrivilege.Resource, requiredPrivilege.Verb)
			}

			methodPrivileges[info.FullMethod] = requiredPrivilegesMask
		}

		if userPrivilegesMask&requiredPrivilegesMask != requiredPrivilegesMask {
			return nil, errstat.MakeGRPCStatus(codes.PermissionDenied, "not enough privileges").Err()
		}

		mexUser := Claims{
			TenantId:   userWithRoles.TenantId,
			AppId:      userWithRoles.AppId,
			UserId:     userWithRoles.UserId,
			Privileges: userPrivilegesMask,
		}
		ctx = context.WithValue(ctx, constants.ContextKeyUserClaims, &mexUser)
		ctx = context.WithValue(ctx, constants.ContextKeyUserID, userWithRoles.UserId)
		ctx = context.WithValue(ctx, constants.ContextKeyTenantID, userWithRoles.TenantId)

		return handler(ctx, req)
	}
}

func getMethodAnnotation[T any](methodName string, extType protoreflect.ExtensionType) (*T, error) {
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(strings.ReplaceAll(methodName[1:], "/", ".")))
	if err != nil {
		// fmt.Printf("ERR: %s\n", err.Error())
		return nil, err
	}

	method, ok := desc.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("no method descriptor")
	}

	options, ok := method.Options().(*descriptorpb.MethodOptions)
	if !ok {
		return nil, fmt.Errorf("no method options")
	}

	if !proto.HasExtension(options, extType) {
		return nil, fmt.Errorf("no extension: %v", extType.TypeDescriptor().FullName())
	}

	t, ok := proto.GetExtension(options, extType).(T)
	if !ok {
		return nil, fmt.Errorf("no extension")
	}

	return &t, nil
}
