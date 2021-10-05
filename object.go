package graphqlgo

import (
	"fmt"
	"reflect"
	"time"

	"github.com/graphql-go/graphql"
)

type CustomNameObject interface {
	Name() string
}

var ObjStorage = make(map[reflect.Type]*graphql.Object)

func GetObjectFromStruct(instance interface{}) (obj *graphql.Object) {
	var t reflect.Type
	if t = reflect.TypeOf(instance); t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if _, ok := ObjStorage[t]; !ok {
		fields := make(graphql.Fields)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			name := field.Name
			tags := getTags(&field)
			if customName, ok := tags["name"]; ok {
				name = customName
			}
			var getFieldType func(field reflect.Type) graphql.Output
			getFieldType = func(field reflect.Type) (fieldType graphql.Output) {
				switch field.Kind() {
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
					fieldType = graphql.Int
				case reflect.Float64, reflect.Float32:
					fieldType = graphql.Float
				case reflect.String:
					fieldType = graphql.String
				case reflect.Bool:
					fieldType = graphql.Boolean
				case reflect.Slice:
					fieldType = graphql.NewList(getFieldType(field.Elem()))
				case reflect.Struct, reflect.Ptr:
					if field.Kind() == reflect.Ptr {
						field = field.Elem()
					}
					if field == reflect.TypeOf(time.Time{}) {
						fieldType = graphql.DateTime
					} else {
						if _, ok := ObjStorage[field]; !ok {
							ObjStorage[field] = GetObjectFromStruct(reflect.New(field).Interface())
						}
						fieldType = ObjStorage[field]
					}
				default:
					panic(fmt.Errorf("cannot recognize type of field %s in object definition %s", field.Name(), t.Name()))
				}
				if _, ok := tags["id"]; ok {
					fieldType = graphql.ID
				}
				if _, ok := tags["not null"]; ok {
					fieldType = graphql.NewNonNull(fieldType)
				}
				return fieldType
			}
			fields[name] = &graphql.Field{
				Name: name,
				Type: getFieldType(field.Type),
			}
		}
		name := t.Name()
		if customNameObject, ok := instance.(CustomNameObject); ok {
			name = customNameObject.Name()
		}
		ObjStorage[t] = graphql.NewObject(graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		})

	}
	return ObjStorage[t]
}
