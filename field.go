package graphqlgo

import (
	"fmt"
	"reflect"
	"time"

	"github.com/graphql-go/graphql"
)

func getFieldType(fieldtype reflect.Type, tags map[string]string, cat typeCategory) (t graphql.Type) {
	fieldtype = unwrapPtr(fieldtype)
	if explicitType, ok := tags["type"]; ok {
		switch explicitType {
		case "id":
			t = graphql.ID
		case "string":
			t = graphql.String
		case "int":
			t = graphql.Int
		case "float":
			t = graphql.Float
		case "bool", "boolean":
			t = graphql.Boolean
		case "date":
			t = graphql.DateTime
		default:
			panic(fmt.Errorf("type %s not recognized", explicitType))
		}
	} else {
		switch fieldtype.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			t = graphql.Int
		case reflect.Float64, reflect.Float32:
			t = graphql.Float
		case reflect.String:
			t = graphql.String
		case reflect.Bool:
			t = graphql.Boolean
		case reflect.Slice:
			t = graphql.NewList(getFieldType(fieldtype.Elem(), make(map[string]string), cat))
		case reflect.Struct, reflect.Ptr:
			if fieldtype == reflect.TypeOf(time.Time{}) {
				t = graphql.DateTime
			} else {
				t = getObjectType(fieldtype, cat)
			}
		default:
			panic(fmt.Errorf("cannot recognize type of field %s", fieldtype.Name()))
		}
	}
	return t
}
