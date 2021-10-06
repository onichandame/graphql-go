package core

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

type Object struct {
	Named
	Describable
	graphql.Type
	Fields []*Field
}

func GetObjectFromsStruct(instance interface{}, opts ...interface{}) (obj Object) {
	structType := unwrapPtr(reflect.TypeOf(instance))
	obj.Named.Name = structType.Name()
	for i, v := range opts {
		switch i {
		case 0:
			if v.(string) != "" {
				obj.Named.Name = v.(string)
			}
		case 1:
			obj.Describable.Description = v.(string)
		}
	}
	return obj
}
