package graphqlgo

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	type Object struct {
		Str          string
		StrRenamed   string `graphqlgo:"name=str"`
		StrDescribed string `graphqlgo:"description=str described"`
	}
	t.Run("input", func(t *testing.T) {
		object := getObjectType(reflect.TypeOf(&Object{}), INPUT)
		assert.NotNil(t, object)
		assert.IsType(t, &graphql.InputObject{}, object)
		inputobject := object.(*graphql.InputObject)
		assert.Equal(t, "Object", inputobject.Name())
		assert.Contains(t, inputobject.Fields(), "Str")
		assert.Equal(t, graphql.String, inputobject.Fields()["Str"].Type)
		assert.Contains(t, inputobject.Fields(), "str")
		assert.Equal(t, graphql.String, inputobject.Fields()["StrDescribed"].Type)
		assert.Equal(t, "str described", inputobject.Fields()["StrDescribed"].Description())
	})
	t.Run("output", func(t *testing.T) {
		object := getObjectType(reflect.TypeOf(&Object{}), OUTPUT)
		assert.NotNil(t, object)
		assert.IsType(t, &graphql.Object{}, object)
		outputobject := object.(*graphql.Object)
		assert.Equal(t, "Object", outputobject.Name())
		assert.Contains(t, outputobject.Fields(), "Str")
		assert.Equal(t, graphql.String, outputobject.Fields()["Str"].Type)
		assert.Contains(t, outputobject.Fields(), "str")
		assert.Equal(t, graphql.String, outputobject.Fields()["StrDescribed"].Type)
		assert.Equal(t, "str described", outputobject.Fields()["StrDescribed"].Description)
	})
}
