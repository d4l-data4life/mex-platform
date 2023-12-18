package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"google.golang.org/grpc/codes"

	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/telemetry"
	"github.com/d4l-data4life/mex/mex/shared/utils"

	pbAuth "github.com/d4l-data4life/mex/mex/services/auth/endpoints/auth/pb"
)

type ServiceOptions struct {
	ServiceTag string
	Log        L.Logger

	Redis *redis.Client

	ClientID        string
	ProducerGroupID string
	ConsumerGroupID string

	SigningKeyPEMReader io.Reader
	SignatureAlg        string
	ClientSecrets       []string
	RedirectURIs        []string
	GrantFlows          []string

	AuthCodeValidity     time.Duration
	AccessTokenValidity  time.Duration
	RefreshTokenValidity time.Duration

	TelemetryService *telemetry.Service
}

type Service struct {
	options    ServiceOptions
	signingKey jwk.Key

	grantFlowClientCredentialsEnabled bool
	grantFlowAuthorizationCodeEnabled bool
	grantFlowRefreshTokenEnabled      bool

	pbAuth.UnimplementedAuthServer
}

type BILog struct {
	GrantType string `json:"grant-type"`
}

func NewService(options ServiceOptions) (*Service, error) {
	if options.Log == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	svc := Service{options: options}

	// This key pair will be used for signing and validating JWTs that this service issues.
	pem, err := io.ReadAll(options.SigningKeyPEMReader)
	if err != nil {
		return nil, err
	}

	signingKey, err := parseSigningKey(pem, options.SignatureAlg)
	if err != nil {
		return nil, err
	}

	svc.signingKey = signingKey

	if len(options.GrantFlows) == 0 {
		options.Log.Warn(context.Background(), L.Message("no grant flows specified; OAuth endpoints will not work"))
	}

	for _, flow := range options.GrantFlows {
		switch flow {
		case grantAuthorizationCode:
			svc.grantFlowAuthorizationCodeEnabled = true
		case grantClientCredentials:
			svc.grantFlowClientCredentialsEnabled = true
		case grantRefreshToken:
			svc.grantFlowRefreshTokenEnabled = true
		default:
			options.Log.Warn(context.Background(), L.Messagef("unknown grant flow: '%s'", flow))
		}
	}

	if len(options.RedirectURIs) == 0 {
		options.Log.Warn(context.Background(), L.Message("no redirect URIs specified; authorization code grant flow will now work"))
	}

	return &svc, nil
}

func (svc *Service) Authorize(ctx context.Context, request *pbAuth.AuthorizeRequest) (*pbAuth.AuthorizeResponse, error) {
	if !svc.grantFlowAuthorizationCodeEnabled {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "authorization code grant flow not enabled").Err()
	}

	if request.ClientId != svc.options.ClientID {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "incorrect client ID").Err()
	}

	if request.ResponseType != responseTypeAuthCode {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, fmt.Sprintf("unsupported response code: %s", request.ResponseType)).Err()
	}

	if request.CodeChallengeMethod != codeChallengeMethod {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, fmt.Sprintf("unsupported challenge method code: %s", request.CodeChallengeMethod)).Err()
	}

	if request.ResponseMode == "" {
		request.ResponseMode = responseModeQuery // default value
	}

	if request.ResponseMode != responseModeQuery && request.ResponseMode != responseModeFragment {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, fmt.Sprintf("unsupported response mode: %s", request.ResponseMode)).Err()
	}

	if request.CodeChallenge == "" {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "challenge missing").Err()
	}

	if !utils.Contains(svc.options.RedirectURIs, request.RedirectUri) {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "redirect URI not whitelisted").Err()
	}

	authCode := randomString()

	hashName := redisAuthHashName(request.CodeChallenge)
	cmdHSet := svc.options.Redis.HSet(ctx, hashName,
		"state", request.State,
		"response_mode", request.ResponseMode,
		"redirect_uri", request.RedirectUri,
		"scope", request.Scope,
		"authCode", authCode,
		"client_id", request.ClientId,
	)
	if cmdHSet.Err() != nil {
		return nil, E.MakeGRPCStatus(codes.Internal, "could not store authentication data in Redis", E.Cause(cmdHSet.Err())).Err()
	}

	cmdExpire := svc.options.Redis.Expire(ctx, hashName, svc.options.AuthCodeValidity)
	if cmdExpire.Err() != nil {
		return nil, E.MakeGRPCStatus(codes.Internal, "could set expiration on authentication data in Redis", E.Cause(cmdExpire.Err())).Err()
	}

	hints.HintHTTPLocation(ctx, fmt.Sprintf("%s%scode=%s&state=%s", request.RedirectUri, fragmentOrQuery(request.ResponseMode), authCode, request.State))
	hints.HintHTTPStatusCode(ctx, http.StatusFound) // 302

	svc.options.Log.BIEvent(ctx, L.BIActivity("auth-authorize"), L.BIData(BILog{GrantType: "authorization_code"}))
	return &pbAuth.AuthorizeResponse{}, nil
}

func (svc *Service) Token(ctx context.Context, request *pbAuth.TokenRequest) (*pbAuth.TokenResponse, error) {
	switch request.GrantType {
	case grantAuthorizationCode:
		return svc.tokenViaAuthorizationCodeFlow(ctx, request)
	case grantClientCredentials:
		return svc.tokenViaClientCredentialsFlow(ctx, request)
	case grantRefreshToken:
		return svc.tokenViaRefreshTokenFlow(ctx, request)
	default:
		return nil, E.MakeGRPCStatus(
			codes.Unauthenticated,
			"unsupported grant type",
			E.DevMessagef("unsupported grant type: '%s'", request.GrantType),
		).Err()
	}
}

func (svc *Service) tokenViaAuthorizationCodeFlow(ctx context.Context, request *pbAuth.TokenRequest) (*pbAuth.TokenResponse, error) {
	if !svc.grantFlowAuthorizationCodeEnabled {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "authorization code grant flow not enabled").Err()
	}

	if request.ClientId != svc.options.ClientID {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "incorrect client ID").Err()
	}

	if request.ClientSecret != "" {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "client secret was given; not allowed in authorization code flow").Err()
	}

	if request.Code == "" {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "authorization code not specified").Err()
	}

	if request.CodeVerifier == "" {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "challenge verifier not specified").Err()
	}

	challenge := sha256B64Safe(request.CodeVerifier)

	cmdHGetAll := svc.options.Redis.HGetAll(ctx, redisAuthHashName(challenge))
	if cmdHGetAll.Err() != nil {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "invalid verifier").Err()
	}

	redisHashValues := cmdHGetAll.Val()

	if authCode, ok := redisHashValues["authCode"]; ok {
		if authCode != request.Code {
			return nil, E.MakeGRPCStatus(codes.Unauthenticated, "incorrect authorization code").Err()
		}
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "authorization code not found in authorization data").Err()
	}

	if clientID, ok := redisHashValues["client_id"]; ok {
		if clientID != request.ClientId {
			return nil, E.MakeGRPCStatus(codes.Unauthenticated, "client ID mismatch").Err()
		}
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "client ID not found in authorization data").Err()
	}

	if redirectURI, ok := redisHashValues["redirect_uri"]; ok {
		if redirectURI != request.RedirectUri {
			return nil, E.MakeGRPCStatus(codes.Unauthenticated, "redirect URI mismatch").Err()
		}
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "redirect URI not found in authorization data").Err()
	}

	accessTokenData := TokenData{
		OID:       "anonymous",
		GrantFlow: grantTypeAuthorizationCode,
		AppID:     svc.options.ClientID,
		Groups:    []string{svc.options.ConsumerGroupID},
		Email:     "anonymous",
		Subject:   "anonymous",
	}

	signedAccessTokenString, err := accessTokenData.toJWT(svc.signingKey, svc.options.AccessTokenValidity)
	if err != nil {
		return nil, err
	}

	refreshTokenString := randomString()
	err = accessTokenData.toRedis(ctx, svc.options.Redis, refreshTokenString, svc.options.RefreshTokenValidity)
	if err != nil {
		return nil, err
	}

	svc.options.Log.BIEvent(ctx, L.BIActivity("token"), L.BIData(struct{ GrantType string }{GrantType: "authorization_code"}))

	hints.HintHTTPCacheControl(ctx, "no-store")
	hints.HintHTTPPragma(ctx, "no-cache")

	return &pbAuth.TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  signedAccessTokenString,
		ExpiresIn:    uint32(svc.options.AccessTokenValidity.Seconds()),
		RefreshToken: refreshTokenString,
	}, nil
}

func (svc *Service) tokenViaClientCredentialsFlow(ctx context.Context, request *pbAuth.TokenRequest) (*pbAuth.TokenResponse, error) {
	if !svc.grantFlowClientCredentialsEnabled {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "client credentials grant flow not enabled").Err()
	}

	if request.ClientId != svc.options.ClientID {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "incorrect client ID").Err()
	}

	if !utils.Contains(svc.options.ClientSecrets, request.ClientSecret) {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "incorrect client secret").Err()
	}

	accessTokenData := TokenData{
		OID:       "anonymous",
		GrantFlow: grantTypeClientCredentials,
		AppID:     svc.options.ClientID,
		Groups:    []string{svc.options.ConsumerGroupID, svc.options.ProducerGroupID},
		Subject:   "anonymous",
	}

	signedAccessTokenString, err := accessTokenData.toJWT(svc.signingKey, svc.options.AccessTokenValidity)
	if err != nil {
		return nil, err
	}

	refreshTokenString := randomString()
	err = accessTokenData.toRedis(ctx, svc.options.Redis, refreshTokenString, svc.options.RefreshTokenValidity)
	if err != nil {
		return nil, err
	}

	svc.options.Log.BIEvent(ctx, L.BIActivity("token"), L.BIData(BILog{GrantType: "client_credentials"}))

	hints.HintHTTPCacheControl(ctx, "no-store")
	hints.HintHTTPPragma(ctx, "no-cache")

	return &pbAuth.TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  signedAccessTokenString,
		ExpiresIn:    uint32(svc.options.AccessTokenValidity.Seconds()),
		RefreshToken: refreshTokenString,
	}, nil
}

func (svc *Service) tokenViaRefreshTokenFlow(ctx context.Context, request *pbAuth.TokenRequest) (*pbAuth.TokenResponse, error) {
	if !svc.grantFlowRefreshTokenEnabled {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "refresh token grant flow not enabled").Err()
	}

	if request.ClientId != svc.options.ClientID {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "incorrect client ID").Err()
	}

	if !utils.Contains(svc.options.ClientSecrets, request.ClientSecret) {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "incorrect client secret").Err()
	}

	if request.RefreshToken == "" {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token specified").Err()
	}

	accessTokenData, err := TokenDataFromRedis(ctx, svc.options.Redis, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	signedAccessTokenString, err := accessTokenData.toJWT(svc.signingKey, svc.options.AccessTokenValidity)
	if err != nil {
		return nil, err
	}

	hints.HintHTTPCacheControl(ctx, "no-store")
	hints.HintHTTPPragma(ctx, "no-cache")

	return &pbAuth.TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  signedAccessTokenString,
		ExpiresIn:    uint32(svc.options.AccessTokenValidity.Seconds()),
		RefreshToken: request.RefreshToken, // reuse refresh token
	}, nil
}

// This method makes the Service an rdb.TopicSubscriber
func (svc *Service) Message(ctx context.Context, topic string, configHash string) {
	svc.options.Log.Info(ctx, L.Messagef("message: %s: %s", topic, configHash))

	// The auth service does not really need anything from the config, so we just respond with the hash in green.
	svc.options.TelemetryService.SetStatus(statuspb.Color_GREEN, configHash)
}
