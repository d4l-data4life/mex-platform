package cfg

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Test_renderValue_NoSecret(t *testing.T) {
	require.Equal(t, "int(-1234)", renderValue(reflect.ValueOf(int(-1234)), false))
	require.Equal(t, "int8(-12)", renderValue(reflect.ValueOf(int8(-12)), false))
	require.Equal(t, "int16(-1234)", renderValue(reflect.ValueOf(int16(-1234)), false))
	require.Equal(t, "int32(-1234)", renderValue(reflect.ValueOf(int32(-1234)), false))
	require.Equal(t, "int64(-1234)", renderValue(reflect.ValueOf(int64(-1234)), false))

	require.Equal(t, "uint(1234)", renderValue(reflect.ValueOf(uint(1234)), false))
	require.Equal(t, "'{'", renderValue(reflect.ValueOf(uint8(123)), false))
	require.Equal(t, "uint16(1234)", renderValue(reflect.ValueOf(uint16(1234)), false))
	require.Equal(t, "uint32(1234)", renderValue(reflect.ValueOf(uint32(1234)), false))
	require.Equal(t, "uint64(1234)", renderValue(reflect.ValueOf(uint64(1234)), false))

	require.Equal(t, "float32(12.340000)", renderValue(reflect.ValueOf(float32(12.34)), false))
	require.Equal(t, "float64(12.340000)", renderValue(reflect.ValueOf(float64(12.34)), false))

	require.Equal(t, "bool(true)", renderValue(reflect.ValueOf(true), false))

	require.Equal(t, "string(\"hello\")", renderValue(reflect.ValueOf("hello"), false))

	require.Equal(t, "[4]uint8('T', 'e', 's', 't')", renderValue(reflect.ValueOf([]byte("Test")), false))
	require.Equal(t, "[4]uint64(uint64(1), uint64(2), uint64(3), uint64(4))", renderValue(reflect.ValueOf([]uint64{1, 2, 3, 4}), false))

	require.Equal(t, "Duration(2m3s)", renderValue(reflect.ValueOf(durationpb.New(123*time.Second)), false))

	foo := "test"
	require.Equal(t, "*string(test)", renderValue(reflect.ValueOf(&foo), false))

	require.Equal(t, "struct({0 0 0 <nil> <nil>})", renderValue(reflect.ValueOf(Abc{}), false))
}

func Test_renderValue_Secret(t *testing.T) {
	require.Equal(t, "int(redacted)", renderValue(reflect.ValueOf(int(-1234)), true))
	require.Equal(t, "int32(redacted)", renderValue(reflect.ValueOf(int32(-1234)), true))
	require.Equal(t, "int64(redacted)", renderValue(reflect.ValueOf(int64(-1234)), true))

	require.Equal(t, "uint(redacted)", renderValue(reflect.ValueOf(uint(1234)), true))
	require.Equal(t, "uint32(redacted)", renderValue(reflect.ValueOf(uint32(1234)), true))
	require.Equal(t, "uint64(redacted)", renderValue(reflect.ValueOf(uint64(1234)), true))

	require.Equal(t, "float32(redacted)", renderValue(reflect.ValueOf(float32(12.34)), true))
	require.Equal(t, "float64(redacted)", renderValue(reflect.ValueOf(float64(12.34)), true))

	require.Equal(t, "bool(redacted)", renderValue(reflect.ValueOf(true), true))

	require.Equal(t, "string(redacted 5 chars)", renderValue(reflect.ValueOf("hello"), true))

	require.Equal(t, "[4]uint8(redacted)", renderValue(reflect.ValueOf([]byte("Test")), true))
	require.Equal(t, "[4]uint64(redacted)", renderValue(reflect.ValueOf([]uint64{1, 2, 3, 4}), true))

	require.Equal(t, "Duration(redacted)", renderValue(reflect.ValueOf(durationpb.New(123*time.Second)), true))

	foo := "test"
	require.Equal(t, "*string(redacted)", renderValue(reflect.ValueOf(&foo), true))

	require.Equal(t, "struct(redacted)", renderValue(reflect.ValueOf(Abc{}), true))
}
