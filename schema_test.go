package graphqlgo_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	graphqlgo "github.com/onichandame/graphql-go"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `graphqlgo:"name=name"`
	Org  *Org   `graphqlgo:"name=org"`
}

var users = make([]User, 0)

type Org struct {
}
type UserInput struct {
	Name string `graphqlgo:"name=name,not null"`
}

func (*User) Queryusers(_ graphql.ResolveParams) []User {
	return users
}

func (*User) MutationcreateOneUser(params graphql.ResolveParams, _ *UserInput) User {
	return User{
		Name: params.Args["name"].(string),
	}
}

func TestSchema(t *testing.T) {
	schema := graphqlgo.GetSchema(graphqlgo.SchemaConfig{
		Resolvers: []interface{}{&User{}},
	})
	assert.NotNil(t, schema)
}
