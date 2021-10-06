package core

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

type Field struct {
}

func GetField(args struct {
	Field   *reflect.StructField
	Resolve graphql.FieldResolveFn
}) (field graphql.Field) {
	tags := getTags(args.Field)
	if name, ok := tags["name"]; ok {
		field.Name = name
	} else {
		field.Name = args.Field.Name
	}
	if description, ok := tags["description"]; ok {
		field.Description = description
	}
	field.Resolve = args.Resolve
	return field
}
