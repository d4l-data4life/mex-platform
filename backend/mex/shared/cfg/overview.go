package cfg

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/d4l-data4life/mex/mex/shared/known/configpb"
)

const base64Suffix = "_B64"

type ConfigOverview interface {
	FilterByEffTag(tag string) ConfigOverview
	Len() int
	Get(i int) FieldOverview
}

type FieldOverview struct {
	ProtoName    protoreflect.FullName
	Kind         protoreflect.Kind
	EnumFullName protoreflect.FullName

	GoPath string
	GoName string

	GenEnvName   string
	AltEnvName   string
	EffEnvName   string
	AliasEnvName string

	Default  string
	Secret   bool
	Repeated bool

	FldTags []string
	MsgTags []string
	EffTags []string

	Title       string
	Summary     string
	Description string

	IgnoreK8s bool
}

type internalOverview []FieldOverview

var durationMessageName = (&durationpb.Duration{}).ProtoReflect().Descriptor().FullName()

func GetMessageOverview(msgDesc protoreflect.MessageDescriptor, globalPrefix string) (ConfigOverview, error) {
	mainMessage := proto.GetExtension(msgDesc.ParentFile().Options().(*descriptorpb.FileOptions), configpb.E_MainMessage).(string)
	if mainMessage == "" {
		return nil, fmt.Errorf("file has no main message option")
	}

	if mainMessage != string(msgDesc.Name()) {
		return nil, fmt.Errorf("message %s (%s) is not marked as main message (but %s is)", msgDesc.Name(), msgDesc.FullName(), mainMessage)
	}

	overview := internalOverview{}

	for i := 0; i < msgDesc.Fields().Len(); i++ {
		overview = append(overview, getFieldOverview(msgDesc.Fields().Get(i), globalPrefix, "", "")...)
	}

	return overview, nil
}

func getFieldOverview(fd protoreflect.FieldDescriptor, globalPrefix string, envNameFragment string, goPathFragment string) internalOverview {
	if fd == nil {
		panic("field descriptor is nil")
	}

	if ignoreField(fd) {
		return internalOverview{}
	}

	protoFieldName := string(fd.Name())                               // proto_style_field_name
	envNameFragment = upperSnakeCase(envNameFragment, protoFieldName) // ENV_VAR_STYLE_NAME_WITHOUT_GLOBAL_PREFIX
	goName := strings.ToUpper(fd.JSONName()[:1]) + fd.JSONName()[1:]
	goPathFragment = goPathFragment + "." + goName

	if fd.Kind() == protoreflect.MessageKind && fd.Message().FullName() != durationMessageName {
		msgFields := fd.Message().Fields()
		overview := internalOverview{}
		for i := 0; i < msgFields.Len(); i++ {
			overview = append(overview, getFieldOverview(msgFields.Get(i), globalPrefix, envNameFragment, goPathFragment)...)
		}
		return overview
	}

	fov := FieldOverview{
		ProtoName:   fd.FullName(),
		GoPath:      goPathFragment,
		GoName:      goName,
		Title:       getTitle(fd),
		Summary:     getSummary(fd),
		Description: strings.Join(getDescription(fd), "\n"),
		Secret:      isSecretField(fd),
		Repeated:    fd.Cardinality() == protoreflect.Repeated,
		Kind:        fd.Kind(),
		IgnoreK8s:   ignoreFieldK8s(fd),
	}

	if fov.Kind == protoreflect.EnumKind {
		fov.EnumFullName = fd.Enum().FullName()
	}

	msgOptions := fd.ContainingMessage().Options().(*descriptorpb.MessageOptions)
	fov.MsgTags = proto.GetExtension(msgOptions, configpb.E_Mtags).([]string)
	if fov.MsgTags == nil {
		fov.MsgTags = []string{}
	}

	fldOptions := fd.Options().(*descriptorpb.FieldOptions)
	fov.FldTags = proto.GetExtension(fldOptions, configpb.E_Tags).([]string)
	if fov.FldTags == nil {
		fov.FldTags = []string{}
	}

	if len(fov.FldTags) > 0 {
		fov.EffTags = fov.FldTags
	} else {
		fov.EffTags = fov.MsgTags
	}

	fov.AliasEnvName = getK8sSource(fd)
	if fov.Kind == protoreflect.BytesKind {
		fov.GenEnvName = upperSnakeCase(globalPrefix, envNameFragment+base64Suffix)
		if fov.AliasEnvName == "" {
			fov.AliasEnvName = upperSnakeCase(globalPrefix, envNameFragment)
		}
	} else {
		fov.GenEnvName = upperSnakeCase(globalPrefix, envNameFragment)
	}

	altEnvName := getAlternateEnvVarName(fd)
	if altEnvName != "" {
		fov.AltEnvName = upperSnakeCase(globalPrefix, altEnvName)
		fov.EffEnvName = fov.AltEnvName
	} else {
		fov.EffEnvName = fov.GenEnvName
	}

	fov.Default = getDefault(fd)
	if fov.Default != "" && fov.Secret {
		panic(fmt.Sprintf("a secret cannot have a default: %s", fd.FullName()))
	}

	return internalOverview{fov}
}

func (ov internalOverview) FilterByEffTag(tag string) ConfigOverview {
	if tag == "*" {
		return ov
	}

	subset := internalOverview{}

	for _, fov := range ov {
		if slices.Contains(fov.EffTags, tag) || slices.Contains(fov.EffTags, "*") {
			subset = append(subset, fov)
		}
	}

	return subset
}

func (ov internalOverview) Len() int {
	return len(ov)
}

func (ov internalOverview) Get(i int) FieldOverview {
	return ov[i]
}
