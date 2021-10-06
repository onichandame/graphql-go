package core

import "github.com/graphql-go/graphql"

type Operation struct {
	Named
	Describable
	Inputs  []*Input
	Output  *Object
	Handler interface{}
}
type Query struct {
	Operation
}
type Mutation struct {
	Operation
}

type Resolver struct {
	Named
	Describable
	Queries   []*Query
	Mutations []*Mutation
}

func (r *Resolver) getFields() (queries graphql.Fields, mutations graphql.Fields) {
	queries = make(graphql.Fields)
	mutations = make(graphql.Fields)
	for _, query := range r.Queries {
		args := make(graphql.FieldConfigArgument)
		for _, input := range query.Inputs {
			args[input.Object.Named.Name] = &graphql.ArgumentConfig{
				Type:        input.Type,
				Description: input.Describable.Description,
			}
		}
		queries[query.Operation.Named.Name] = &graphql.Field{
			Description: query.Description,
			Args:        args,
			Type:        query.Output.Type,
		}
	}
	for _, mutation := range r.Mutations {
		args := make(graphql.FieldConfigArgument)
		for _, input := range mutation.Inputs {
			args[input.Object.Named.Name] = &graphql.ArgumentConfig{
				Type:        input.Type,
				Description: input.Describable.Description,
			}
		}
		mutations[mutation.Operation.Named.Name] = &graphql.Field{
			Description: mutation.Description,
			Args:        args,
			Type:        mutation.Output.Type,
		}
	}
	return queries, mutations
}
