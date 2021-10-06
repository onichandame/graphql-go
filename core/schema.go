package core

import "github.com/graphql-go/graphql"

func GetSchemaFromResolvers(resolvers []*Resolver) *graphql.Schema {
	queries := make(graphql.Fields)
	mutations := make(graphql.Fields)
	for _, resolver := range resolvers {
		q, m := resolver.getFields()
		for k, v := range q {
			queries[k] = v
		}
		for k, v := range m {
			mutations[k] = v
		}
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
	}); err == nil {
		return &schema
	} else {
		panic(err)
	}
}
