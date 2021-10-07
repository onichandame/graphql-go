package graphqlgo

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

type Input struct {
	Name string `graphqlgo:"name=name,not null"`
}
type Output struct {
	Greeting string `graphqlgo:"name=greeting"`
}
type Resolver struct{}

func (*Resolver) Querygreet(input graphql.ResolveParams, _ *Input) Output {
	name := input.Args["name"]
	return Output{
		Greeting: fmt.Sprintf("hello %s", name),
	}
}
func TestResolver(t *testing.T) {
	query, _, _ := getResolver(reflect.TypeOf(&Resolver{}))
	assert.NotNil(t, query)
	assert.Contains(t, query, "greet")
	greet := query["greet"]
	greetoutput := getObjectType(reflect.TypeOf(&Output{}), OUTPUT)
	assert.Equal(t, greetoutput, greet.Type)
}
