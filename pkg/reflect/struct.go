package reflect

import (
	"reflect"

	"github.com/fatih/structtag"
	"github.com/palantir/stacktrace"
)

type StructT struct {
	reflect.Type
	Fields map[string]reflect.StructField
	Tags   map[string]*structtag.Tags
}

func Struct(s any) (*StructT, error) {
	t := reflect.ValueOf(s)
	if t.Kind() == reflect.Pointer {
		return Struct(t.Elem().Interface())
	}
	if t.Kind() != reflect.Struct {
		return nil, stacktrace.NewError("unexpected %T in reflect.Struct call", s)
	}

	fields := make(map[string]reflect.StructField)
	tags := make(map[string]*structtag.Tags)
	exportedFields := reflect.VisibleFields(t.Type())
	for _, field := range exportedFields {
		fields[field.Name] = field
		parsedTags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return nil, stacktrace.Propagate(err, "unable to parse struct '%s' field", field.Name)
		}
		tags[field.Name] = parsedTags
	}

	return &StructT{
		Type:   t.Type(),
		Fields: fields,
		Tags:   tags,
	}, nil
}
