package graphqlgo

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestArgs(t *testing.T) {
	type Args struct {
		Str          string
		StrRenamed   string `graphqlgo:"name=str"`
		StrDescribed string `graphqlgo:"description=str described"`
	}
	args := getArgs(reflect.TypeOf(&Args{}))
	assert.NotNil(t, args)
	assert.Contains(t, args, "Str")
	assert.Equal(t, graphql.String, args["Str"].Type)
	assert.Contains(t, args, "str")
	assert.Contains(t, args, "StrDescribed")
	assert.Equal(t, "str described", args["StrDescribed"].Description)
}
