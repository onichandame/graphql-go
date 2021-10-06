package graphqlgo

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

func getArgs(argstype reflect.Type) (t graphql.FieldConfigArgument) {
	t = make(graphql.FieldConfigArgument)
	for i := 0; i < argstype.NumField(); i++ {
		fieldtype := argstype.Field(i)
		tags := getTags(&fieldtype)
		name := fieldtype.Name
		var description string
		if explicitname, ok := tags["name"]; ok {
			name = explicitname
		}
		if desc, ok := tags["description"]; ok {
			description = desc
		}
		t[name] = &graphql.ArgumentConfig{
			Type:        getFieldType(fieldtype.Type, tags, INPUT),
			Description: description,
		}
	}
	return t
}
