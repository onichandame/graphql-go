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
type Resolver struct{}
type GreetInput struct {
	Name string `graphqlgo:"name=name,not null"`
}
type GreetOutput struct {
	Greeting string `graphqlgo:"name=greeting,not null"`
}
type HandShakeInput struct {
	Hand string `graphqlgo:"name=hand,not null"`
}
type HandShakeOutput struct {
	Hand string `graphqlgo:"name=hand,not null"`
}

func (*Resolver) Querygreet(input graphql.ResolveParams, _ *GreetInput) GreetOutput {
	name := input.Args["name"]
	return GreetOutput{
		Greeting: fmt.Sprintf("hello %s", name),
	}
}

func (*Resolver) Mutationhandshake(input graphql.ResolveParams, _ *HandShakeInput) HandShakeOutput {
	return HandShakeOutput{
		Hand: input.Args["hand"].(string),
	}
}
func TestResolver(t *testing.T) {
	query, mutation, _ := getResolver(reflect.TypeOf(&Resolver{}))
	assert.NotNil(t, query)
	t.Run("query", func(t *testing.T) {
		assert.Contains(t, query, "greet")
		greet := query["greet"]
		greetoutput := getObjectType(reflect.TypeOf(GreetOutput{}), OUTPUT)
		assert.Equal(t, greetoutput, greet.Type)
		msg, _ := greet.Resolve(graphql.ResolveParams{Args: map[string]interface{}{"name": "world"}})
		assert.IsType(t, GreetOutput{}, msg)
		assert.Equal(t, "hello world", msg.(GreetOutput).Greeting)
	})
	t.Run("mutation", func(t *testing.T) {
		assert.Contains(t, mutation, "handshake")
		handshake := mutation["handshake"]
		handshakeoutput := getObjectType(reflect.TypeOf(HandShakeOutput{}), OUTPUT)
		assert.Equal(t, handshakeoutput, handshake.Type)
		hand, _ := handshake.Resolve(graphql.ResolveParams{Args: map[string]interface{}{"hand": "left"}})
		assert.IsType(t, HandShakeOutput{}, hand)
		assert.Equal(t, "left", hand.(HandShakeOutput).Hand)
	})
}
