package cfg

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"

	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

//nolint:nonamedreturns
func InitConfig(log L.Logger, envs EnvVars, globalPrefix string, tag string, obj proto.Message) (dump string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in config parser: %v", r)
		}
	}()

	log.Info(context.Background(), L.Message("parsing: configuration"), L.Phase("config:"+tag))

	overview, err := GetMessageOverview(obj.ProtoReflect().Descriptor(), globalPrefix)
	if err != nil {
		panic(err)
	}

	log.Info(context.Background(), L.Messagef("# of total cfg fields: %d", overview.Len()))
	overview = overview.FilterByEffTag(tag)
	log.Info(context.Background(), L.Messagef("# of %s fields: %d", tag, overview.Len()))

	objR := reflect.ValueOf(obj)
	resetStructPointerFields(objR)

	sb := strings.Builder{}
	for i := 0; i < overview.Len(); i++ {
		ov := overview.Get(i)
		setValue(log, &ov, envs.GetValue(ov.EffEnvName), tag, &sb, objR)
	}

	log.Info(context.Background(), L.Message("parsed: configuration"), L.Phase("config:"+tag))

	dump = sb.String()
	return dump, nil
}

func setValue(log L.Logger, fov *FieldOverview, envVarValue string, tag string, sb *strings.Builder, objR reflect.Value) {
	// Determine whether this field needs setting
	if !utils.Contains(fov.EffTags, tag) && !utils.Contains(fov.EffTags, "*") {
		// Field does not need setting according to annotations, but we issue
		// a warning if an env value is present.
		if envVarValue != "" {
			log.Warn(context.Background(), L.Messagef("variable %s is set, but ignored for service tag '%s'", fov.EffEnvName, tag), L.Phase("config:"+tag))
		}
		return
	}

	effectiveValue := getOrDefault(envVarValue, fov, tag)
	fieldR := reflectGetFieldByPath(objR, fov.GoPath)

	// Set the value in the config message structure.
	if fov.Repeated {
		fieldR.Set(reflect.ValueOf(parseValues(fov, effectiveValue)))
	} else {
		fieldR.Set(reflect.ValueOf(parseValue(fov, effectiveValue)))
	}

	// Produce a human-readable version.
	sb.WriteString(fmt.Sprintf("%-52s: %-30v (%s)\n",
		strings.TrimPrefix(string(fov.ProtoName), "d4l.mex.cfg.Mex"),
		renderValue(fieldR, fov.Secret),
		fov.EffEnvName))
}
