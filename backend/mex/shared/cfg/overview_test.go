package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"

	cfgtestpb "github.com/d4l-data4life/mex/mex/shared/cfg/test/pb"
)

func TestOverview1(t *testing.T) {
	c := cfgtestpb.NonMain{}

	_, err := GetMessageOverview(c.ProtoReflect().Descriptor(), "FOO")
	require.Error(t, err)
}

func TestOverview2(t *testing.T) {
	c := cfgtestpb.TestConfig{}

	overview, err := GetMessageOverview(c.ProtoReflect().Descriptor(), "FOO")
	require.NoError(t, err)

	require.Equal(t, 14, overview.Len())

	tt := []FieldOverview{
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.pi"),
			Kind:      protoreflect.DoubleKind,

			GoPath: ".Pi",
			GoName: "Pi",

			GenEnvName: "FOO_PI",
			EffEnvName: "FOO_PI",

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Database.user"),
			Kind:      protoreflect.StringKind,

			GoPath: ".Db.User",
			GoName: "User",

			GenEnvName: "FOO_DB_USER",
			EffEnvName: "FOO_DB_USER",

			Default: "hjones",

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Database.password"),
			Kind:      protoreflect.StringKind,

			GoPath: ".Db.Password",
			GoName: "Password",

			GenEnvName: "FOO_DB_PASSWORD",
			EffEnvName: "FOO_DB_PASSWORD",

			Secret: true,

			FldTags: []string{"AAA", "BBB"},
			MsgTags: []string{"*"},
			EffTags: []string{"AAA", "BBB"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Database.port"),
			Kind:      protoreflect.Uint32Kind,

			GoPath: ".Db.Port",
			GoName: "Port",

			GenEnvName: "FOO_DB_PORT",
			AltEnvName: "FOO_PORT_NUMBER",
			EffEnvName: "FOO_PORT_NUMBER",

			Default: "5432",

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Database.search_path"),
			Kind:      protoreflect.StringKind,

			GoPath: ".Db.SearchPath",
			GoName: "SearchPath",

			GenEnvName: "FOO_DB_SEARCH_PATH",
			EffEnvName: "FOO_DB_SEARCH_PATH",

			Default:  "mex,public",
			Repeated: true,

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Database.use_ssl"),
			Kind:      protoreflect.BoolKind,

			GoPath: ".Db.UseSsl",
			GoName: "UseSsl",

			GenEnvName: "FOO_DB_USE_SSL",
			EffEnvName: "FOO_DB_USE_SSL",

			Default: "true",

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},

			Title: "Whether to use TLS/SSL",
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Database.timeout"),
			Kind:      protoreflect.MessageKind,

			GoPath: ".Db.Timeout",
			GoName: "Timeout",

			GenEnvName: "FOO_DB_TIMEOUT",
			EffEnvName: "FOO_DB_TIMEOUT",

			Default: "2s",

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},

			IgnoreK8s: true,
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Server.timeout"),
			Kind:      protoreflect.MessageKind,

			GoPath: ".Server.Timeout",
			GoName: "Timeout",

			GenEnvName:   "FOO_SERVER_TIMEOUT",
			EffEnvName:   "FOO_SERVER_TIMEOUT",
			AliasEnvName: "GLOBAL_TIMEOUT",

			Default: "5s",

			FldTags: []string{"bar", "wom"},
			MsgTags: []string{"foo", "bar"},
			EffTags: []string{"bar", "wom"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Server.max_header_bytes"),
			Kind:      protoreflect.Fixed32Kind,

			GoPath: ".Server.MaxHeaderBytes",
			GoName: "MaxHeaderBytes",

			GenEnvName: "FOO_SERVER_MAX_HEADER_BYTES",
			EffEnvName: "FOO_SERVER_MAX_HEADER_BYTES",

			Default: "2097152",

			FldTags: []string{},
			MsgTags: []string{"foo", "bar"},
			EffTags: []string{"foo", "bar"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.Server.signing_private_key_pem"),
			Kind:      protoreflect.BytesKind,

			GoPath: ".Server.SigningPrivateKeyPem",
			GoName: "SigningPrivateKeyPem",

			GenEnvName:   "FOO_SERVER_SIGNING_PRIVATE_KEY_PEM_B64",
			EffEnvName:   "FOO_SERVER_SIGNING_PRIVATE_KEY_PEM_B64",
			AliasEnvName: "FOO_SERVER_SIGNING_PRIVATE_KEY_PEM",

			Secret: true,

			FldTags: []string{},
			MsgTags: []string{"foo", "bar"},
			EffTags: []string{"foo", "bar"},

			Summary:     "summary",
			Description: "desc 1\ndesc 2\ndesc 3",
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.public_keys"),
			Kind:      protoreflect.BytesKind,

			GoPath: ".PublicKeys",
			GoName: "PublicKeys",

			GenEnvName:   "FOO_PUBLIC_KEYS_B64",
			EffEnvName:   "FOO_PUBLIC_KEYS_B64",
			AliasEnvName: "FOO_PUBLIC_KEYS",

			Repeated: true,

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.secret_keys"),
			Kind:      protoreflect.BytesKind,

			GoPath: ".SecretKeys",
			GoName: "SecretKeys",

			GenEnvName:   "FOO_SECRET_KEYS_B64",
			EffEnvName:   "FOO_SECRET_KEYS_B64",
			AliasEnvName: "FOO_SECRET_KEYS",

			Secret:   true,
			Repeated: true,

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.fibonacci"),
			Kind:      protoreflect.Sint64Kind,

			GoPath: ".Fibonacci",
			GoName: "Fibonacci",

			GenEnvName: "FOO_FIBONACCI",
			EffEnvName: "FOO_FIBONACCI",

			Repeated: true,

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
		{
			ProtoName: protoreflect.FullName("d4l.mex.cfg.test.TestConfig.constants"),
			Kind:      protoreflect.DoubleKind,

			GoPath: ".Constants",
			GoName: "Constants",

			GenEnvName: "FOO_CONSTANTS",
			EffEnvName: "FOO_CONSTANTS",

			Repeated: true,

			FldTags: []string{},
			MsgTags: []string{"*"},
			EffTags: []string{"*"},
		},
	}

	for i, fov := range tt {
		require.Equal(t, fov, overview.Get(i))
	}
}
