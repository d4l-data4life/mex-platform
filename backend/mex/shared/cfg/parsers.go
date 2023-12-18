package cfg

import (
	b64 "encoding/base64"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/d4l-data4life/mex/mex/shared/utils"
)

const (
	listSep  = ","
	emptySet = "∅"
)

func parseValue(fov *FieldOverview, rawValue string) any {
	switch fov.Kind {
	case protoreflect.StringKind:
		return rawValue

	case protoreflect.BoolKind:
		b, err := strconv.ParseBool(rawValue)
		if err != nil {
			panic(err)
		}
		return b

	case protoreflect.Int32Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		if i < math.MinInt32 || i > math.MaxInt32 {
			panic(fmt.Sprintf("outside int32 bounds: %d", i))
		}
		//#nosec
		return int32(i)

	case protoreflect.Int64Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		return int64(i)

	case protoreflect.Sint32Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		if i < math.MinInt32 || i > math.MaxInt32 {
			panic(fmt.Sprintf("outside int32 bounds: %d", i))
		}
		//#nosec
		return int32(i)

	case protoreflect.Sint64Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		return int64(i)

	case protoreflect.Uint32Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		if i < 0 || i > math.MaxUint32 {
			panic(fmt.Sprintf("outside uint32 bounds: %d", i))
		}
		//#nosec
		return uint32(i)

	case protoreflect.Uint64Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		return uint64(i)

	case protoreflect.Fixed32Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		if i < 0 || i > math.MaxUint32 {
			panic(fmt.Sprintf("outside uint32 bounds: %d", i))
		}
		//#nosec
		return uint32(i)

	case protoreflect.Fixed64Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		return uint64(i)

	case protoreflect.Sfixed32Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		if i < math.MinInt32 || i > math.MaxInt32 {
			panic(fmt.Sprintf("outside int32 bounds: %d", i))
		}
		//#nosec
		return int32(i)

	case protoreflect.Sfixed64Kind:
		i, err := strconv.Atoi(rawValue)
		if err != nil {
			panic(fmt.Sprintf("%s: error: %s", fov.ProtoName, err.Error()))
		}
		return int64(i)

	case protoreflect.FloatKind:
		f, err := strconv.ParseFloat(rawValue, 32)
		if err != nil {
			panic(err)
		}
		return float32(f)

	case protoreflect.DoubleKind:
		f, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			panic(err)
		}
		return f

	case protoreflect.BytesKind:
		buf, err := b64.StdEncoding.DecodeString(rawValue)
		if err != nil {
			panic(err)
		}
		return buf

	case protoreflect.MessageKind:
		d, err := time.ParseDuration(rawValue)
		if err != nil {
			panic(err)
		}

		dur := durationpb.New(d)
		return dur

	case protoreflect.EnumKind:
		enumType, err := protoregistry.GlobalTypes.FindEnumByName(fov.EnumFullName)
		if err != nil {
			panic(err)
		}

		valueDesc := enumType.Descriptor().Values().ByName(protoreflect.Name(rawValue))
		if valueDesc == nil {
			panic(fmt.Sprintf("enum %s: no value found for: %s", string(fov.ProtoName), rawValue))
		}

		return enumType.New(valueDesc.Number())

	default:
		panic(fmt.Sprintf("unsupported kind: %v", fov.Kind))
	}
}

func parseValues(fov *FieldOverview, rawValue string) any {
	// We need a non-empty string to indicate an empty list because we want to
	// make sure an empty env var (or an empty default value) cause an error.
	// However, for a list, the empty string would be a valid representation.
	rawValue = TrimEmpty(rawValue)

	parts := strings.Split(rawValue, listSep)
	switch fov.Kind {
	case protoreflect.StringKind:
		return utils.Map(parts, func(s string) string { return parseValue(fov, s).(string) })
	case protoreflect.BoolKind:
		return utils.Map(parts, func(s string) bool { return parseValue(fov, s).(bool) })
	case protoreflect.Int32Kind:
		return utils.Map(parts, func(s string) int32 { return parseValue(fov, s).(int32) })
	case protoreflect.Int64Kind:
		return utils.Map(parts, func(s string) int64 { return parseValue(fov, s).(int64) })
	case protoreflect.Sint32Kind:
		return utils.Map(parts, func(s string) int32 { return parseValue(fov, s).(int32) })
	case protoreflect.Sint64Kind:
		return utils.Map(parts, func(s string) int64 { return parseValue(fov, s).(int64) })
	case protoreflect.Uint32Kind:
		return utils.Map(parts, func(s string) uint32 { return parseValue(fov, s).(uint32) })
	case protoreflect.Uint64Kind:
		return utils.Map(parts, func(s string) uint64 { return parseValue(fov, s).(uint64) })
	case protoreflect.Fixed32Kind:
		return utils.Map(parts, func(s string) uint32 { return parseValue(fov, s).(uint32) })
	case protoreflect.Fixed64Kind:
		return utils.Map(parts, func(s string) uint64 { return parseValue(fov, s).(uint64) })
	case protoreflect.Sfixed32Kind:
		return utils.Map(parts, func(s string) int32 { return parseValue(fov, s).(int32) })
	case protoreflect.Sfixed64Kind:
		return utils.Map(parts, func(s string) int64 { return parseValue(fov, s).(int64) })
	case protoreflect.FloatKind:
		return utils.Map(parts, func(s string) float32 { return parseValue(fov, s).(float32) })
	case protoreflect.DoubleKind:
		return utils.Map(parts, func(s string) float64 { return parseValue(fov, s).(float64) })
	case protoreflect.BytesKind:
		return utils.Map(parts, func(s string) []byte { return parseValue(fov, s).([]byte) })
	case protoreflect.EnumKind:
		return utils.Map(parts, func(s string) protoreflect.EnumNumber { return parseValue(fov, s).(protoreflect.EnumNumber) })
	default:
		panic("unsupported kind: " + fmt.Sprint(fov.Kind))
	}
}

func renderValue(val reflect.Value, secret bool) string {
	if secret {
		switch val.Kind() {
		case reflect.String:
			return fmt.Sprintf("string(redacted %d chars)", len(val.String()))
		case reflect.Slice:
			return fmt.Sprintf("[%d]%s(redacted)", val.Len(), val.Type().Elem().Kind())
		case reflect.Ptr:
			if val.Elem().Type().String() == "durationpb.Duration" {
				return "Duration(redacted)"
			}
			return fmt.Sprintf("*%s(redacted)", val.Elem().Kind())

		default:
			return fmt.Sprintf("%s(redacted)", val.Kind())
		}
	}

	switch val.Kind() {
	case reflect.String:
		return fmt.Sprintf("string(\"%s\")", val.String())
	case reflect.Bool:
		return fmt.Sprintf("bool(%v)", val.Bool())

	case reflect.Int:
		return fmt.Sprintf("int(%d)", val.Int())
	case reflect.Int8:
		return fmt.Sprintf("int8(%d)", val.Int())
	case reflect.Int16:
		return fmt.Sprintf("int16(%d)", val.Int())
	case reflect.Int32:
		return fmt.Sprintf("int32(%d)", val.Int())
	case reflect.Int64:
		return fmt.Sprintf("int64(%d)", val.Int())

	case reflect.Uint:
		return fmt.Sprintf("uint(%d)", val.Uint())
	case reflect.Uint8:
		return fmt.Sprintf("'%s'", string([]uint8{uint8(val.Uint())}))
	case reflect.Uint16:
		return fmt.Sprintf("uint16(%d)", val.Uint())
	case reflect.Uint32:
		return fmt.Sprintf("uint32(%d)", val.Uint())
	case reflect.Uint64:
		return fmt.Sprintf("uint64(%d)", val.Uint())

	case reflect.Float32:
		return fmt.Sprintf("float32(%f)", val.Float())
	case reflect.Float64:
		return fmt.Sprintf("float64(%f)", val.Float())

	case reflect.Slice:
		return fmt.Sprintf("[%d]%s(%s)", val.Len(), val.Type().Elem().Kind(), renderSlice(val))

	case reflect.Ptr:
		if val.Elem().Type().String() == "durationpb.Duration" {
			return fmt.Sprintf("Duration(%s)", val.Interface().(*durationpb.Duration).AsDuration().String())
		}
		return fmt.Sprintf("*%s(%v)", val.Elem().Kind(), val.Elem())

	default:
		return fmt.Sprintf("%s(%v)", val.Kind(), val)
	}
}

func renderSlice(val reflect.Value) string {
	parts := make([]string, val.Len())

	for i := 0; i < val.Len(); i++ {
		parts[i] = renderValue(val.Index(i), false)
	}

	return strings.Join(parts, ", ")
}
