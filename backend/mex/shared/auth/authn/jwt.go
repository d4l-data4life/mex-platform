package authn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/d4l-data4life/mex/mex/shared/auth"
	"github.com/d4l-data4life/mex/mex/shared/keys"
	"github.com/d4l-data4life/mex/mex/shared/known/securitypb"
)

type jwtAuthenticator struct {
	validator keys.TokenValidator
	clientID  string

	consumerGroupID string
	producerGroupID string
}

type tokenHeader struct {
	tokenType  string
	tokenValue string
}

func NewJWTAuthenticator(validator keys.TokenValidator, clientID string, consumerGroupID string, producerGroupID string) (auth.RequestAuthenticator, error) {
	return &jwtAuthenticator{
		validator:       validator,
		clientID:        clientID,
		consumerGroupID: consumerGroupID,
		producerGroupID: producerGroupID,
	}, nil
}

func (a *jwtAuthenticator) Authenticate(ctx context.Context, request any) (*securitypb.UserWithRoles, error) {
	authHeader, err := extractAuthorizationHeader(ctx)
	if err != nil {
		return nil, err
	}

	if authHeader.tokenType != "bearer" {
		return nil, status.Error(codes.Unauthenticated, "authorization header does not contain a Bearer token")
	}

	token, err := a.validator.ValidateJWT(authHeader.tokenValue)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	user := securitypb.UserWithRoles{}

	if value, ok := (*token).Get("appid"); ok {
		user.AppId = value.(string)
		if user.AppId != a.clientID {
			return nil, status.Error(codes.Unauthenticated, "incorrect client/app ID")
		}
	} else {
		return nil, status.Error(codes.Unauthenticated, "missing claim: 'appid'")
	}

	if value, ok := (*token).Get("oid"); ok {
		user.UserId = value.(string)
	} else {
		return nil, status.Error(codes.Unauthenticated, "missing claim: 'oid'")
	}

	// Distinguish technical vs human user
	var errAppidacr, errAzpacr error
	roles, errAppidacr := determineRoles(a.consumerGroupID, a.producerGroupID, token, "appidacr")
	if errAppidacr != nil {
		roles, errAzpacr = determineRoles(a.consumerGroupID, a.producerGroupID, token, "azpacr")
		if errAzpacr != nil {
			return nil, status.Error(
				codes.Unauthenticated,
				fmt.Sprintf("error when reading claim 'appidacr' (%s); error when reading 'azpacr' claim (%s)", errAppidacr.Error(), errAzpacr.Error()),
			)
		}
	}

	user.Roles = roles
	return &user, nil
}

func determineRoles(consumerGroupID, producerGroupID string, token *jwt.Token, claim string) ([]string, error) {
	if acr, ok := (*token).Get(claim); ok {
		// See for values: https://docs.microsoft.com/en-us/azure/active-directory/develop/access-tokens
		switch acr {
		case "0": // human user token, check groups
			if rawGroups, ok := (*token).Get("groups"); ok {
				if groups, ok := rawGroups.([]interface{}); ok {
					roles := []string{}
					for _, g := range groups {
						if g == consumerGroupID {
							roles = append(roles, auth.RoleConsumer)
						}
						if g == producerGroupID {
							roles = append(roles, auth.RoleProducer)
						}
					}
					return roles, nil
				}
				return nil, errors.New("human user, but 'groups' is of wrong type")
			}
			return nil, errors.New("human user, but 'groups' claim is missing")

		case "1": // technical user token
			return []string{auth.RoleConsumer, auth.RoleProducer}, nil
		default:
			return nil, fmt.Errorf("unsupported '%s' value", claim)
		}
	} else {
		return nil, fmt.Errorf("claim missing: '%s'", claim)
	}
}

func extractAuthorizationHeader(ctx context.Context) (*tokenHeader, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "incoming gRPC metadata missing")
	}

	authorizations := md.Get("authorization")
	if len(authorizations) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header missing")
	}

	authorization := authorizations[0]
	parts := strings.Split(authorization, " ")
	if len(parts) != 2 {
		return nil, status.Error(codes.Unauthenticated, "authorization header malformed")
	}

	return &tokenHeader{
		tokenType:  strings.ToLower(parts[0]),
		tokenValue: parts[1],
	}, nil
}
