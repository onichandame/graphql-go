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
	greetoutput := getObjectType(reflect.TypeOf(Output{}), OUTPUT)
	assert.Equal(t, greetoutput, greet.Type)
	msg, _ := greet.Resolve(graphql.ResolveParams{Args: map[string]interface{}{"name": "world"}})
	assert.IsType(t, Output{}, msg)
	assert.Equal(t, "hello world", msg.(Output).Greeting)
}
