package graphqlgo

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
)

var CustomOutObjects = make(map[reflect.Type]*graphql.Object)
var CustomInObjects = make(map[reflect.Type]*graphql.InputObject)

type typeCategory int

const (
	INPUT  typeCategory = iota
	OUTPUT typeCategory = iota
)

func getObjectType(objtype reflect.Type, cat typeCategory) (t graphql.Type) {
	objtype = unwrapPtr(objtype)
	var ok bool
	switch cat {
	case INPUT:
		t, ok = CustomInObjects[objtype]
	case OUTPUT:
		t, ok = CustomOutObjects[objtype]
	default:
		panic(fmt.Errorf("object category %d not valid", cat))
	}
	if !ok {
		var fields interface{}
		// parse fields
		switch cat {
		case INPUT:
			// TODO: get default value from method
			inputfields := make(graphql.InputObjectConfigFieldMap)
			for i := 0; i < objtype.NumField(); i++ {
				fieldtype := objtype.Field(i)
				tags := getTags(&fieldtype)
				name := fieldtype.Name
				var description string
				if explicitname, ok := tags["name"]; ok {
					name = explicitname
				}
				if desc, ok := tags["description"]; ok {
					description = desc
				}
				inputfields[name] = &graphql.InputObjectFieldConfig{
					Type:        getFieldType(fieldtype.Type, tags, cat),
					Description: description,
				}
			}
			fields = inputfields
		case OUTPUT:
			// TODO: get subscribe/resolve from method/tag
			outputfields := graphql.Fields{}
			for i := 0; i < objtype.NumField(); i++ {
				fieldtype := objtype.Field(i)
				tags := getTags(&fieldtype)
				name := fieldtype.Name
				var description string
				if explicitname, ok := tags["name"]; ok {
					name = explicitname
				}
				if desc, ok := tags["description"]; ok {
					description = desc
				}
				outputfields[name] = &graphql.Field{
					Type:        getFieldType(fieldtype.Type, tags, cat),
					Description: description,
					Name:        name,
				}
			}
			fields = outputfields
		}
		// store type
		switch cat {
		case INPUT:
			CustomInObjects[objtype] = graphql.NewInputObject(graphql.InputObjectConfig{
				Fields: fields,
				Name:   objtype.Name(),
			})
			t = CustomInObjects[objtype]
		case OUTPUT:
			CustomOutObjects[objtype] = graphql.NewObject(graphql.ObjectConfig{
				Fields: fields,
				Name:   objtype.Name(),
			})
			t = CustomOutObjects[objtype]
		}
	}
	return t
}
