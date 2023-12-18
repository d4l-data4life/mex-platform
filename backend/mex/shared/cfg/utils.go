package cfg

import (
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/d4l-data4life/mex/mex/shared/known/configpb"
)

func getOrDefault(value string, fov *FieldOverview, tag string) string {
	if value != "" {
		return value
	}

	if fov.Default != "" {
		return fov.Default
	}

	panic(fmt.Sprintf("field '%s' and its default value are empty, but relevant for service '%s'", fov.ProtoName, tag))
}

func getDefault(fd protoreflect.FieldDescriptor) string {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_Opts) {
		x := proto.GetExtension(options, configpb.E_Opts).(*configpb.Options)
		return x.Default
	}

	return ""
}

func getK8sSource(fd protoreflect.FieldDescriptor) string {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_K8S) {
		x := proto.GetExtension(options, configpb.E_K8S).(*configpb.Export)
		return x.Source
	}

	return ""
}

func getSummary(fd protoreflect.FieldDescriptor) string {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_Desc) {
		x := proto.GetExtension(options, configpb.E_Desc).(*configpb.Descriptor)
		return x.Summary
	}

	return ""
}

func getTitle(fd protoreflect.FieldDescriptor) string {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_Desc) {
		x := proto.GetExtension(options, configpb.E_Desc).(*configpb.Descriptor)
		return x.Title
	}

	return ""
}

func getDescription(fd protoreflect.FieldDescriptor) []string {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_Desc) {
		x := proto.GetExtension(options, configpb.E_Desc).(*configpb.Descriptor)
		return x.Description
	}

	return []string{}
}

func getAlternateEnvVarName(fd protoreflect.FieldDescriptor) string {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_Opts) {
		x := proto.GetExtension(options, configpb.E_Opts).(*configpb.Options)
		return x.Env
	}
	return ""
}

func ignoreField(fd protoreflect.FieldDescriptor) bool {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_Opts) {
		x := proto.GetExtension(options, configpb.E_Opts).(*configpb.Options)
		return x.Ignore
	}
	return false
}

func ignoreFieldK8s(fd protoreflect.FieldDescriptor) bool {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_K8S) {
		x := proto.GetExtension(options, configpb.E_K8S).(*configpb.Export)
		return x.Ignore
	}
	return false
}

func isSecretField(fd protoreflect.FieldDescriptor) bool {
	options := fd.Options().(*descriptorpb.FieldOptions)
	if proto.HasExtension(options, configpb.E_Opts) {
		x := proto.GetExtension(options, configpb.E_Opts).(*configpb.Options)
		return x.Secret
	}
	return false
}

// The function name should really be UPPER_SNAKE_CASE ;)
func upperSnakeCase(prefix string, rest string) string {
	if prefix == "" {
		return strings.ToUpper(rest)
	}
	return strings.ToUpper(prefix + "_" + rest)
}

// This function takes a (reflection value of) pointer to a struct and
// recursively descents into it setting all pointer fields to a pointer
// toa newly created respective type instance.
// It leaves all other fields untouched.
func resetStructPointerFields(objR reflect.Value) {
	if objR.Kind() == reflect.Ptr {
		if objR.IsNil() {
			zeroR := reflect.New(objR.Type().Elem())
			objR.Set(zeroR)
		}

		e := objR.Elem()
		if e.Kind() != reflect.Struct {
			return
		}

		for i := 0; i < e.Type().NumField(); i++ {
			fd := e.Type().Field(i)
			resetStructPointerFields(e.FieldByName(fd.Name))
		}
	}
}

func reflectGetFieldByPath(objR reflect.Value, path string) reflect.Value {
	if path == "" {
		panic("path is empty")
	}

	if path[0:1] == "." {
		path = path[1:]
	}

	parts := strings.Split(path, ".")
	fieldR := objR

	for _, part := range parts {
		if fieldR.Kind() == reflect.Ptr {
			fieldR = fieldR.Elem()
		}

		fieldR = fieldR.FieldByName(part)
	}

	return fieldR
}

func TrimEmpty(s string) string {
	return strings.Trim(s, " "+emptySet)
}

// A string is considered empty if it becomes the empty string after
// trimming spaces and the empty set symbol ∅.
func StringIsEmpty(s string) bool {
	return TrimEmpty(s) == ""
}

func BytesAreEmpty(b []byte) bool {
	if b == nil {
		return true
	}
	return StringIsEmpty(string(b))
}
