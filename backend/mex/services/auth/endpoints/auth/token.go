package auth

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"google.golang.org/grpc/codes"

	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/uuid"
)

type TokenData struct {
	OID       string
	GrantFlow string
	AppID     string
	Groups    []string
	Scopes    []string
	Subject   string
	Email     string
}

func (t *TokenData) toJWT(signingKey jwk.Key, lifetime time.Duration) (string, error) {
	b := jwt.NewBuilder().
		JwtID(uuid.MustNewV4()).
		Claim("oid", t.OID).
		Claim("appidacr", t.GrantFlow).
		Claim("appid", t.AppID).
		Claim("groups", t.Groups).
		Subject(t.Subject).
		Expiration(time.Now().Add(lifetime))

	if t.Email != "" {
		b = b.Claim("email", t.Email)
	}
	if len(t.Scopes) > 0 {
		b = b.Claim("scope", t.Scopes) // singular 'scope' according to spec
	}

	j, err := b.Build()
	if err != nil {
		return "", err
	}

	signedJWT, err := jwt.Sign(j, jwt.WithKey(jwa.RS256, signingKey))
	if err != nil {
		return "", err
	}

	return string(signedJWT), nil
}

func (t *TokenData) toRedis(ctx context.Context, rdb *redis.Client, key string, lifetime time.Duration) error {
	hashName := redisAuthRefreshTokenHashName(key)
	cmdHSet := rdb.HSet(ctx, hashName,
		"oid", t.OID,
		"appidacr", t.GrantFlow,
		"appid", t.AppID,
		"groups", strings.Join(t.Groups, " "),
		"scopes", strings.Join(t.Scopes, " "),
		"sub", t.Subject,
		"email", t.Email,
	)
	if cmdHSet.Err() != nil {
		return E.MakeGRPCStatus(codes.Internal, "could not store authentication data in Redis", E.Cause(cmdHSet.Err())).Err()
	}

	cmdExpire := rdb.Expire(ctx, hashName, lifetime)
	if cmdExpire.Err() != nil {
		return E.MakeGRPCStatus(codes.Internal, "could set expiration on authentication data in Redis", E.Cause(cmdExpire.Err())).Err()
	}

	return nil
}

func TokenDataFromRedis(ctx context.Context, rdb *redis.Client, key string) (*TokenData, error) {
	hashName := redisAuthRefreshTokenHashName(key)
	cmdHGetAll := rdb.HGetAll(ctx, hashName)

	if cmdHGetAll.Err() != nil {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "invalid refresh token").Err()
	}

	redisHashValues := cmdHGetAll.Val()

	token := TokenData{}

	if oid, ok := redisHashValues["oid"]; ok {
		token.OID = oid
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token data").Err()
	}

	if appidacr, ok := redisHashValues["appidacr"]; ok {
		token.GrantFlow = appidacr
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token data").Err()
	}

	if appid, ok := redisHashValues["appid"]; ok {
		token.AppID = appid
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token data").Err()
	}

	if groups, ok := redisHashValues["groups"]; ok {
		token.Groups = strings.Split(groups, " ")
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token data").Err()
	}

	if scopes, ok := redisHashValues["scopes"]; ok {
		token.Scopes = strings.Split(scopes, " ")
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token data").Err()
	}

	if sub, ok := redisHashValues["sub"]; ok {
		token.Subject = sub
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token data").Err()
	}

	if email, ok := redisHashValues["email"]; ok {
		token.Email = email
	} else {
		return nil, E.MakeGRPCStatus(codes.Unauthenticated, "no refresh token data").Err()
	}

	return &token, nil
}
