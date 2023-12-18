package cfg

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	cfgtestpb "github.com/d4l-data4life/mex/mex/shared/cfg/test/pb"
	"github.com/d4l-data4life/mex/mex/shared/log"
)

func TestEmptyEnvs(t *testing.T) {
	envs := map[string]string{
		"MEX_PI":                                 "3.1415",
		"MEX_VERSION_BUILD":                      "v1.2.3",
		"MEX_VERSION_DESC":                       "Service",
		"MEX_WEB_METRICS_PATH":                   "/metrics",
		"MEX_SERVER_TIMEOUT":                     "10s",
		"MEX_SERVER_SIGNING_PRIVATE_KEY_PEM_B64": "SGVsbG8=",
		"MEX_SEARCH_PATHS":                       "foo,bar,wom",
		"MEX_FIBONACCI":                          "1,1,2,3,5,8,13,21",
		"MEX_CONSTANTS":                          "3.1415,2.71,1.141",
		"MEX_PUBLIC_KEYS_B64":                    "Zm9v,YmFy,d29t",
		"MEX_SECRET_KEYS_B64":                    "Zm9v,YmFy,d29t",
	}

	var c cfgtestpb.TestConfig

	_, err := InitConfig(&log.NullLogger{}, MockedEnvs(envs), "MEX", "bar", &c)
	require.NoError(t, err)

	require.Equal(t, 3.1415, c.Pi)
	require.Equal(t, 10*time.Second, c.Server.Timeout.AsDuration())
	require.Equal(t, []byte("Hello"), c.Server.SigningPrivateKeyPem)
	require.Equal(t, []string{"mex", "public"}, c.Db.SearchPath)
	require.Equal(t, []int64{1, 1, 2, 3, 5, 8, 13, 21}, c.Fibonacci)
	require.Equal(t, []float64{3.1415, 2.71, 1.141}, c.Constants)
	require.Equal(t, [][]byte{[]byte("foo"), []byte("bar"), []byte("wom")}, c.PublicKeys)
	require.Equal(t, [][]byte{[]byte("foo"), []byte("bar"), []byte("wom")}, c.SecretKeys)
}

func TestNonConfigRoot(t *testing.T) {
	envs := map[string]string{}
	var c cfgtestpb.TestConfig_Database // the field 'password' has no default and we should get an error

	_, err := InitConfig(&log.NullLogger{}, MockedEnvs(envs), "mex", "test", &c)
	require.Error(t, err)
	t.Logf("error: %s", err.Error())
}

type Abc struct {
	A int
	B int
	C int
	D *int
	E *string
}

type User struct {
	Name  string
	ID    uint64
	Alpha *Abc
}

type Foo struct {
	Bar   int
	Wom   string
	Admin *User
}

// Some reflection refresher.
func TestReflection(t *testing.T) {
	var foo Foo

	r := reflect.ValueOf(&foo)
	require.Equal(t, reflect.Ptr, r.Kind())
	require.Equal(t, "*cfg.Foo", r.Type().String())

	e := r.Elem()
	require.Equal(t, reflect.Struct, e.Kind())

	fBar := e.FieldByName("Bar")
	fBar.SetInt(1234)

	fWom := e.FieldByName("Wom")
	fWom.SetString("Hello")

	fAdmin := e.FieldByName("Admin")
	require.Equal(t, reflect.Ptr, fAdmin.Kind())
	require.Equal(t, "*cfg.User", fAdmin.Type().String())
	require.Equal(t, "cfg.User", fAdmin.Type().Elem().String())
	require.True(t, fAdmin.IsNil())

	objAdmin := reflect.New(fAdmin.Type().Elem())
	require.Equal(t, reflect.Ptr, objAdmin.Kind())
	require.Equal(t, "*cfg.User", objAdmin.Type().String())

	objAdmin.Elem().FieldByName("Name").SetString("hjones")
	objAdmin.Elem().FieldByName("ID").SetUint(4711)

	fAdmin.Set(objAdmin)

	require.Equal(t, 1234, foo.Bar)
	require.Equal(t, "Hello", foo.Wom)
	require.Equal(t, "hjones", foo.Admin.Name)
	require.Equal(t, uint64(4711), foo.Admin.ID)
	require.Nil(t, foo.Admin.Alpha)
}

func TestZeroOut(t *testing.T) {
	// Do not touch values that are not pointers to a struct.
	bar := []string{"Hello", "World"}
	resetStructPointerFields(reflect.ValueOf(&bar))
	require.Len(t, bar, 2)
	require.Equal(t, "Hello", bar[0])
	require.Equal(t, "World", bar[1])

	foo := Foo{Bar: 1234}
	resetStructPointerFields(reflect.ValueOf(&foo))
	require.Equal(t, 1234, foo.Bar)
	require.Equal(t, "", foo.Wom)
	require.Equal(t, "", foo.Admin.Name)
	require.Equal(t, uint64(0), foo.Admin.ID)
	require.Equal(t, 0, foo.Admin.Alpha.A)
	require.Equal(t, 0, foo.Admin.Alpha.B)
	require.Equal(t, 0, foo.Admin.Alpha.C)
	require.Equal(t, 0, *foo.Admin.Alpha.D)
	require.Equal(t, "", *foo.Admin.Alpha.E)

	var msg = &cfgtestpb.TestConfig{}
	resetStructPointerFields(reflect.ValueOf(msg))
	require.Nil(t, msg.Constants)
	require.Equal(t, &durationpb.Duration{}, msg.Server.Timeout)
	require.Equal(t, uint32(0), msg.Server.MaxHeaderBytes)
	require.Nil(t, msg.Server.SigningPrivateKeyPem)
	require.Equal(t, "", msg.Db.Password)
	require.Equal(t, uint32(0), msg.Db.Port)
	require.Equal(t, false, msg.Db.UseSsl)
	require.Equal(t, float64(0), msg.Pi)
	require.Nil(t, msg.SecretKeys)
}
