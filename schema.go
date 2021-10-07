package graphqlgo

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

type SchemaConfig struct {
	Resolvers []interface{}
}

func GetSchema(props SchemaConfig) *graphql.Schema {
	queries := make(graphql.Fields)
	mutations := make(graphql.Fields)
	subscriptions := make(graphql.Fields)
	for _, resolver := range props.Resolvers {
		q, m, s := getResolver(reflect.TypeOf(resolver))
		mergeResolvers(queries, q)
		mergeResolvers(mutations, m)
		mergeResolvers(subscriptions, s)
	}
	if schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queries,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutations,
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Subscription",
			Fields: subscriptions,
		}),
	}); err != nil {
		panic(err)
	} else {
		return &schema
	}
}
