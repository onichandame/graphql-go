package graphqlgo

import (
	"reflect"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	t.Run("input fields", func(t *testing.T) {
		type Element struct{}
		type Input struct {
			Str      string
			Int      int
			Float    float32
			ID       uint `graphqlgo:"type=id"`
			Bool     bool
			Date     *time.Time
			SliceStr []string
			Elements []*Element
		}
		inputtype := reflect.TypeOf(Input{})
		t.Run("plain string", func(t *testing.T) {
			plainfield, _ := inputtype.FieldByName("Str")
			plainfieldtags := getTags(&plainfield)
			field := getFieldType(plainfield.Type, plainfieldtags, INPUT)
			assert.NotNil(t, field)
			assert.Equal(t, graphql.String, field)
		})
		t.Run("int number", func(t *testing.T) {
			intfield, _ := inputtype.FieldByName("Int")
			intfieldtags := getTags(&intfield)
			field := getFieldType(intfield.Type, intfieldtags, INPUT)
			assert.NotNil(t, field)
			assert.Equal(t, graphql.Int, field)
		})
		t.Run("float number", func(t *testing.T) {
			floatfield, _ := inputtype.FieldByName("Float")
			floatfieldtags := getTags(&floatfield)
			field := getFieldType(floatfield.Type, floatfieldtags, INPUT)
			assert.NotNil(t, field)
			assert.Equal(t, graphql.Float, field)
		})
		t.Run("number id", func(t *testing.T) {
			idfield, _ := inputtype.FieldByName("ID")
			idfieldtags := getTags(&idfield)
			field := getFieldType(idfield.Type, idfieldtags, INPUT)
			assert.NotNil(t, field)
			assert.Equal(t, graphql.ID, field)
		})
		t.Run("bool", func(t *testing.T) {
			boolfield, _ := inputtype.FieldByName("Bool")
			boolfieldtags := getTags(&boolfield)
			field := getFieldType(boolfield.Type, boolfieldtags, INPUT)
			assert.NotNil(t, field)
			assert.Equal(t, graphql.Boolean, field)
		})
		t.Run("date", func(t *testing.T) {
			datefield, _ := inputtype.FieldByName("Date")
			datefieldtags := getTags(&datefield)
			field := getFieldType(datefield.Type, datefieldtags, INPUT)
			assert.NotNil(t, field)
			assert.Equal(t, graphql.DateTime, field)
		})
		t.Run("slice string", func(t *testing.T) {
			slicefield, _ := inputtype.FieldByName("SliceStr")
			slicefieldtags := getTags(&slicefield)
			field := getFieldType(slicefield.Type, slicefieldtags, INPUT)
			assert.NotNil(t, field)
			assert.IsType(t, &graphql.List{}, field)
			assert.Equal(t, graphql.String, field.(*graphql.List).OfType)
		})
		t.Run("custom object", func(t *testing.T) {
			objfield, _ := inputtype.FieldByName("Elements")
			objfieldtags := getTags(&objfield)
			field := getFieldType(objfield.Type, objfieldtags, INPUT)
			assert.NotNil(t, field)
			assert.IsType(t, &graphql.List{}, field)
			//assert.Equal(t, graphql.String, field.(*graphql.List).OfType)
		})
	})
}
