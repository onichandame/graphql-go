package graphqlgo_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	graphqlgo "github.com/onichandame/graphql-go"
	"github.com/stretchr/testify/assert"
)

type BaseSimple struct{}
type BaseRenamed struct{}

func (s *BaseRenamed) Name() string {
	return "baseRenamed"
}

type StrSimple struct {
	Str string
}
type StrRenamed struct {
	Str string `gimgraphql:"name=str"`
}
type IntSimple struct {
	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64
}
type FloatSimple struct {
	Float32 float32
	Float64 float64
}
type BoolSimple struct {
	Bool bool
}
type DateSimple struct {
	Date    time.Time
	DatePtr *time.Time
}
type ErrorSimple struct {
	Chan chan (interface{})
}
type SliceSimple struct {
	Sli []string
}
type SliceComposite struct {
	Sli    []StrSimple
	SliPtr []*StrSimple
}
type CompositeSimple struct {
	StrSimple  StrSimple
	StrRenamed *StrRenamed
}
type CompositeComposite struct {
	CompositeSimple    CompositeSimple
	CompositeSimplePtr *CompositeSimple
}

// TODO: refine composite/slice test
func TestObject(t *testing.T) {
	t.Run("can handle struct with default name", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&BaseSimple{})
		assert.Equal(t, "BaseSimple", obj.Name())
	})
	t.Run("can handle struct with custom name", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&BaseRenamed{})
		assert.Equal(t, "baseRenamed", obj.Name())
	})
	t.Run("can handle renamed field", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&StrRenamed{})
		assert.Contains(t, obj.Fields(), "str")
	})
	t.Run("can handle string field", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&StrSimple{})
		assert.Contains(t, obj.Fields(), "Str")
		assert.Equal(t, graphql.String, obj.Fields()["Str"].Type)
	})
	t.Run("can handle int field", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&IntSimple{})
		assert.Len(t, obj.Fields(), 5)
		for _, field := range obj.Fields() {
			assert.Equal(t, graphql.Int, field.Type)
		}
	})
	t.Run("can handle float field", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&FloatSimple{})
		assert.Len(t, obj.Fields(), 2)
		for _, field := range obj.Fields() {
			assert.Equal(t, graphql.Float, field.Type)
		}
	})
	t.Run("can handle bool field", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&BoolSimple{})
		assert.Contains(t, obj.Fields(), "Bool")
		assert.Equal(t, graphql.Boolean, obj.Fields()["Bool"].Type)
	})
	t.Run("can handle date field", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&DateSimple{})
		assert.Contains(t, obj.Fields(), "Date")
		assert.Equal(t, graphql.DateTime, obj.Fields()["Date"].Type)
		assert.Contains(t, obj.Fields(), "DatePtr")
		assert.Equal(t, graphql.DateTime, obj.Fields()["DatePtr"].Type)
	})
	t.Run("can handle field composite of primitives", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&CompositeSimple{})
		assert.Contains(t, obj.Fields(), "StrSimple")
		assert.Contains(t, obj.Fields(), "StrRenamed")
	})
	t.Run("can handle field composite of composites", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&CompositeComposite{})
		assert.Contains(t, obj.Fields(), "CompositeSimple")
		assert.Contains(t, obj.Fields(), "CompositeSimplePtr")
	})
	t.Run("can handle slice of primitives", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&SliceSimple{})
		assert.Contains(t, obj.Fields(), "Sli")
		assert.Equal(t, reflect.TypeOf(obj.Fields()["Sli"].Type), reflect.TypeOf(&graphql.List{}))
	})
	t.Run("can handle slice of composites", func(t *testing.T) {
		obj := graphqlgo.GetObjectFromStruct(&SliceComposite{})
		assert.Contains(t, obj.Fields(), "Sli")
		assert.Equal(t, reflect.TypeOf(obj.Fields()["Sli"].Type), reflect.TypeOf(&graphql.List{}))
		assert.Contains(t, obj.Fields(), "SliPtr")
		assert.Equal(t, reflect.TypeOf(obj.Fields()["SliPtr"].Type), reflect.TypeOf(&graphql.List{}))
	})
	t.Run("can panic in case of not supported field", func(t *testing.T) {
		assert.Panics(t, func() { graphqlgo.GetObjectFromStruct(&ErrorSimple{}) })
	})
}
